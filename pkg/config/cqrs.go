package config

import (
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"github.com/pperaltaisern/financing/internal/esrc/esrcwatermill"
	"github.com/pperaltaisern/financing/internal/watermillzap"
	"github.com/pperaltaisern/financing/pkg/command"
	"github.com/pperaltaisern/financing/pkg/eventhandler"
	"github.com/pperaltaisern/financing/pkg/intevent"
	"go.uber.org/zap"
)

func BuildCqrsFacade(log *zap.Logger, repos Repositories) (*cqrs.Facade, *message.Router, error) {
	amqpConfig := LoadAMQPConfig()
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
		GenerateEventsTopic: generateEventsTopic,
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

func BuildEventBus(log *zap.Logger) (*cqrs.EventBus, error) {
	amqpConfig := LoadAMQPConfig()

	publisher, err := amqpConfig.BuildEventPublisher(log)
	if err != nil {
		return nil, err
	}

	return cqrs.NewEventBus(
		publisher,
		generateEventsTopic,
		buildCommandsEventMarshaler())
}

func buildCommandsEventMarshaler() cqrs.CommandEventMarshaler {
	return esrcwatermill.RelayEventMarshaler{
		CmdMarshaler: cqrs.JSONMarshaler{
			GenerateName: cqrs.StructName,
		},
	}
}

func generateEventsTopic(eventName string) string {
	// because we are using PubSub RabbitMQ config, we can use one topic for all events
	return "events"
}