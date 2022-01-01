package main

import (
	"context"
	"flag"
	"fmt"
	espostgres "ledger/internal/esrc/postgres"
	"ledger/internal/esrc/relay"
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

	pool, err := pgxpool.Connect(context.Background(), config.Postgres.ConnectionString)
	if err != nil {
		log.Fatal("err connecting to Postgres: %v", zap.Error(err))
	}

	repos, err := PostgresRepositories(pool)
	if err != nil {
		log.Fatal("err building Postgres repositories: %v", zap.Error(err))
	}

	cqrsFacade, messageRouter, err := CqrsFacade(config.AMQP, repos, log)
	if err != nil {
		log.Fatal("err building CQRS facade: %v", zap.Error(err))
	}

	m := Main{
		log:           log,
		messageRouter: messageRouter,
		relayer: relay.NewRelayer(
			espostgres.NewEventStoreOutbox(pool),
			Publisher{cqrsFacade.EventBus()},
		),
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

	inv := intevent.InvestorRegistered{
		ID:      financing.NewID(),
		Name:    "INVESTOR_1",
		Balance: 100,
	}
	err = cqrsFacade.EventBus().Publish(context.Background(), inv)
	if err != nil {
		log.Error("err creating investor", zap.Error(err))
	}

	if false {
		iss := intevent.IssuerRegistered{
			ID:   financing.NewID(),
			Name: "ISSUER_1",
		}
		err = cqrsFacade.EventBus().Publish(context.Background(), iss)
		if err != nil {
			log.Error("err creating issuer", zap.Error(err))
		}

		err = cqrsFacade.EventBus().Publish(context.Background(), iss)
		if err != nil {
			log.Error("err creating issuer", zap.Error(err))
		}

		cmd := command.SellInvoice{
			IssuerID:    financing.NewID(),
			AskingPrice: 20,
		}
		err = cqrsFacade.CommandBus().Send(context.Background(), cmd)
		if err != nil {
			log.Error("err SellInvoice", zap.Error(err))
		}
	}

	log.Info("terminated", zap.Error(<-errC))
	m.Close()
	log.Info("closed gracefully")
}

type Publisher struct {
	bus *cqrs.EventBus
}

func (p Publisher) Publish(ctx context.Context, e relay.Event) error {
	fmt.Println("** Event relayed: ", e)
	return p.bus.Publish(ctx, e)
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

func PostgresRepositories(pool *pgxpool.Pool) (Repositories, error) {
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

	eventsPublisher, err := amqp.NewPublisher(amqp.NewDurablePubSubConfig(config.Address, nil), wmlog)
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
	relayer       *relay.Relayer
}

func (m *Main) Run() error {
	go m.relayer.Run()
	return m.messageRouter.Run(context.Background())
}

func (m *Main) Close() {
	m.relayer.Stop()
	err := m.messageRouter.Close()
	if err != nil {
		m.log.Error("err clossing message router %v: err")
	}
	m.log.Sync()
}
