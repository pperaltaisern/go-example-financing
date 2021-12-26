package main

import (
	"context"
	"flag"
	"fmt"
	"ledger/internal/watermillzap"
	"ledger/pkg/command"
	"ledger/pkg/eventhandler"
	"ledger/pkg/financing"
	"ledger/pkg/intevent"
	"ledger/pkg/postgres"
	"os"
	"os/signal"

	"github.com/ThreeDotsLabs/watermill-amqp/pkg/amqp"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"
)

func main() {
	config, log, err := ParseFlags()
	if err != nil {
		panic(err)
	}

	repos, err := PostgresRepositories(config.Postgres)
	if err != nil {
		log.Fatal("err building Postgres repositories: %v", zap.Error(err))
	}

	_, messageRouter, err := CqrsFacade(config.AMQP, repos, log)
	if err != nil {
		log.Fatal("err building CQRS facade: %v", zap.Error(err))
	}

	m := Main{
		log:           log,
		messageRouter: messageRouter,
	}

	errC := make(chan error, 2)
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		errC <- fmt.Errorf("%s", <-c)
	}()

	go func() {
		errC <- m.Run()
	}()

	log.Info("terminated", zap.Error(<-errC))
	m.Close()
	log.Info("closed gracefully")
}

func ParseFlags() (c Config, log *zap.Logger, err error) {
	configFile := flag.String("config", "config_dev.json", "application settings")
	configDirectory := flag.String("configdir", "./", "config directory")
	flag.Parse()

	c, err = NewConfig(*configDirectory, *configFile)
	if err != nil {
		return
	}

	log, err = zap.NewDevelopmentConfig().Build()
	if err != nil {
		return
	}

	return
}

type Repositories struct {
	Issuers   financing.IssuerRepository
	Investors financing.InvestorRepository
	Invoices  financing.InvoiceRepository
}

func PostgresRepositories(config PostgresConfig) (Repositories, error) {
	pool, err := pgxpool.Connect(context.Background(), config.ConnectionString)
	if err != nil {
		return Repositories{}, err
	}

	repos := Repositories{
		Issuers:   postgres.NewIssuerRepository(pool),
		Investors: postgres.NewInvestorRepository(pool),
		Invoices:  postgres.NewInvoiceRepository(pool),
	}
	return repos, nil
}

func CqrsFacade(config AMQPConfig, repos Repositories, log *zap.Logger) (*cqrs.Facade, *message.Router, error) {
	wmlog := watermillzap.NewLogger(log)

	cqrsMarshaler := cqrs.JSONMarshaler{}
	commandsAMQPConfig := amqp.NewDurableQueueConfig(config.Address)

	commandsPublisher, err := amqp.NewPublisher(commandsAMQPConfig, wmlog)
	if err != nil {
		return nil, nil, err
	}

	commandsSubscriber, err := amqp.NewSubscriber(commandsAMQPConfig, wmlog)
	if err != nil {
		return nil, nil, err
	}

	eventsPublisher, err := amqp.NewPublisher(commandsAMQPConfig, wmlog)
	if err != nil {
		return nil, nil, err
	}

	router, err := message.NewRouter(message.RouterConfig{}, wmlog)
	if err != nil {
		return nil, nil, err
	}
	router.AddMiddleware(middleware.Recoverer)

	facade, err := cqrs.NewFacade(cqrs.FacadeConfig{
		GenerateCommandsTopic: func(commandName string) string {
			// we are using queue RabbitMQ config, so we need to have topic per command type
			return commandName
		},
		CommandHandlers: func(cb *cqrs.CommandBus, eb *cqrs.EventBus) []cqrs.CommandHandler {
			return []cqrs.CommandHandler{
				command.NewBidOnInvoiceHandler(repos.Investors),
				command.NewApproveFinancingHandler(repos.Invoices),
				command.NewReverseFinancingHandler(repos.Invoices),
				command.NewSellInvoiceHandler(repos.Issuers, repos.Invoices),
				command.NewCreateInvestorHandler(repos.Investors),
				command.NewCreateIssuerHandler(repos.Issuers),
			}
		},
		CommandsPublisher: commandsPublisher,
		CommandsSubscriberConstructor: func(handlerName string) (message.Subscriber, error) {
			// we can reuse subscriber, because all commands have separated topics
			return commandsSubscriber, nil
		},
		GenerateEventsTopic: func(eventName string) string {
			// because we are using PubSub RabbitMQ config, we can use one topic for all events
			return "events"

			// we can also use topic per event type
			// return eventName
		},
		EventHandlers: func(cb *cqrs.CommandBus, eb *cqrs.EventBus) []cqrs.EventHandler {
			return []cqrs.EventHandler{
				eventhandler.NewBidOnInvoicePlacedHandler(repos.Invoices),
				eventhandler.NewBidOnInvoiceRejectedHandler(repos.Investors),
				eventhandler.NewInvoiceFinancedHandler(repos.Investors),
				eventhandler.NewInvoiceReversedHandler(repos.Investors),
				intevent.NewInvestorRegisteredHandler(cb),
				intevent.NewIssuerRegisteredHandler(cb),
			}
		},
		EventsSubscriberConstructor: func(handlerName string) (message.Subscriber, error) {
			config := amqp.NewDurablePubSubConfig(
				config.Address,
				amqp.GenerateQueueNameTopicNameWithSuffix(handlerName),
			)

			return amqp.NewSubscriber(config, wmlog)
		},
		EventsPublisher:       eventsPublisher,
		Router:                router,
		CommandEventMarshaler: cqrsMarshaler,
		Logger:                wmlog,
	})
	if err != nil {
		router.Close()
		return nil, nil, err
	}

	return facade, router, nil
}

type Main struct {
	log           *zap.Logger
	messageRouter *message.Router
}

func (m *Main) Run() error {
	return m.messageRouter.Run(context.Background())
}

func (m *Main) Close() {
	err := m.messageRouter.Close()
	if err != nil {
		m.log.Error("err clossing message router %v: err")
	}
	m.log.Sync()
}
