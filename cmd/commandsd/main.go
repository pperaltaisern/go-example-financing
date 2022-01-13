package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/pperaltaisern/financing/internal/esrc/esrcpg"
	"github.com/pperaltaisern/financing/internal/esrc/esrcwatermill"
	"github.com/pperaltaisern/financing/internal/esrc/relay"
	"github.com/pperaltaisern/financing/internal/watermillzap"
	"github.com/pperaltaisern/financing/pkg/command"
	"github.com/pperaltaisern/financing/pkg/config"
	"github.com/pperaltaisern/financing/pkg/eventhandler"
	"github.com/pperaltaisern/financing/pkg/financing"
	"github.com/pperaltaisern/financing/pkg/grpc"
	"github.com/pperaltaisern/financing/pkg/intevent"

	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"
)

func main() {
	log, err := config.LoadLoggerConfig().Build()
	if err != nil {
		panic(err)
	}

	pool, err := config.LoadPostgresConfig().Build()
	if err != nil {
		log.Fatal("err connecting to Postgres: %v", zap.Error(err))
	}
	repos, err := PostgresRepositories(pool)
	if err != nil {
		log.Fatal("err building Postgres repositories: %v", zap.Error(err))
	}

	cqrsFacade, messageRouter, err := CqrsFacade(repos, log)
	if err != nil {
		log.Fatal("err building CQRS facade: %v", zap.Error(err))
	}

	serverConfig := config.LoadServerConfig()
	m := Main{
		log:           log,
		messageRouter: messageRouter,
		commandServer: grpc.NewCommandServer(
			serverConfig.Network,
			serverConfig.Port,
			cqrsFacade.CommandBus(),
		),
		relayer: relay.NewRelayer(
			esrcpg.NewEventStoreOutbox(pool),
			esrcwatermill.NewPublisher(cqrsFacade.EventBus()),
			relay.RelayerWithLogger(log),
		),
	}

	errC := make(chan error, 4)
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		errC <- fmt.Errorf("%s", <-c)
	}()

	m.Run(errC)

	if false {
		// inv := intevent.InvestorRegistered{
		// 	ID:      financing.NewID(),
		// 	Name:    "INVESTOR_1",
		// 	Balance: 100,
		// }
		// err = cqrsFacade.EventBus().Publish(context.Background(), inv)
		// if err != nil {
		// 	log.Error("err creating investor", zap.Error(err))
		// }

		// iss := intevent.IssuerRegistered{
		// 	ID:   financing.NewID(),
		// 	Name: "ISSUER_1",
		// }
		// err = cqrsFacade.EventBus().Publish(context.Background(), iss)
		// if err != nil {
		// 	log.Error("err creating issuer", zap.Error(err))
		// }

		// cmd := command.SellInvoice{
		// 	InvoiceID:   financing.NewID(),
		// 	IssuerID:    financing.NewIDFrom("37bca316-3b73-4caf-8230-6e4f287ab2e1"),
		// 	AskingPrice: 20,
		// }
		// err = cqrsFacade.CommandBus().Send(context.Background(), cmd)
		// if err != nil {
		// 	log.Error("err SellInvoice", zap.Error(err))
		// }

		cmd := command.BidOnInvoice{
			InvoiceID:  financing.NewIDFromString("89a332f9-0cd7-4a43-8770-6bf5027ef1e7"),
			InvestorID: financing.NewIDFromString("ca8573f2-203a-4d9e-bd2c-621edf7b9eed"),
			BidAmount:  35,
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

type Repositories struct {
	Issuers   financing.IssuerRepository
	Investors financing.InvestorRepository
	Invoices  financing.InvoiceRepository
}

func PostgresRepositories(pool *pgxpool.Pool) (Repositories, error) {
	es := esrcpg.NewEventStore(pool)
	repos := Repositories{
		Issuers:   financing.NewIssuerRepository(es),
		Investors: financing.NewInvestorRepository(es),
		Invoices:  financing.NewInvoiceRepository(es),
	}
	return repos, nil
}

func CqrsFacade(repos Repositories, log *zap.Logger) (*cqrs.Facade, *message.Router, error) {
	amqpConfig := config.LoadAMQPConfig()
	wmlog := watermillzap.NewLogger(log)

	cqrsMarshaler := esrcwatermill.RelayEventMarshaler{
		CmdMarshaler: cqrs.JSONMarshaler{
			GenerateName: cqrs.StructName,
		},
	}

	commandsPublisher, err := amqpConfig.BuildCommandPublisher(log)
	if err != nil {
		return nil, nil, err
	}

	commandsSubscriber, err := amqpConfig.BuildCommandSubscriber(log)
	if err != nil {
		return nil, nil, err
	}

	eventsPublisher, err := amqpConfig.BuildEventPublisher(log)
	if err != nil {
		return nil, nil, err
	}

	router, err := message.NewRouter(message.RouterConfig{}, watermillzap.NewLogger(log))
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
			return amqpConfig.BuildEventSubscriber(log, handlerName)
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
	commandServer *grpc.CommandServer
}

func (m *Main) Run(errC chan<- error) {
	go m.relayer.Run()
	go func() { errC <- m.messageRouter.Run(context.Background()) }()
	go func() { errC <- m.commandServer.Open() }()
}

func (m *Main) Close() {
	m.relayer.Stop()
	err := m.messageRouter.Close()
	if err != nil {
		m.log.Error("err clossing message router %v: err")
	}
	m.commandServer.Close()
	m.log.Sync()
}
