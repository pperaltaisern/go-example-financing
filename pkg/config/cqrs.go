package config

import (
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"github.com/pperaltaisern/financing/internal/esrc"
	"github.com/pperaltaisern/financing/internal/esrc/esrcwatermill"
	"github.com/pperaltaisern/financing/internal/watermillzap"
	"github.com/pperaltaisern/financing/pkg/command"
	"github.com/pperaltaisern/financing/pkg/eventhandler"
	"github.com/pperaltaisern/financing/pkg/financing"
	"github.com/pperaltaisern/financing/pkg/intevent"
	"go.uber.org/zap"
)

func BuildCqrsFacade(log *zap.Logger, repos Repositories) (*cqrs.Facade, *message.Router, error) {
	amqpConfig := LoadAMQPConfig()
	wmlog := watermillzap.NewLogger(log)

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
	router.AddMiddleware(errorHandlingMiddleware(log))

	facade, err := cqrs.NewFacade(cqrs.FacadeConfig{
		GenerateCommandsTopic: func(commandName string) string {
			// we are using queue RabbitMQ config, so we need to have topic per command type
			return commandName
		},
		CommandHandlers: func(cb *cqrs.CommandBus, eb *cqrs.EventBus) []cqrs.CommandHandler {
			return []cqrs.CommandHandler{
				esrcwatermill.NewHandler[command.BidOnInvoice](
					command.NewBidOnInvoiceHandler(repos.Investors)),

				esrcwatermill.NewHandler[command.ApproveFinancing](
					command.NewApproveFinancingHandler(repos.Invoices)),

				esrcwatermill.NewHandler[command.ReverseFinancing](
					command.NewReverseFinancingHandler(repos.Invoices)),

				esrcwatermill.NewHandler[command.SellInvoice](
					command.NewSellInvoiceHandler(repos.Issuers, repos.Invoices)),

				esrcwatermill.NewHandler[command.CreateInvestor](
					command.NewCreateInvestorHandler(repos.Investors)),

				esrcwatermill.NewHandler[command.CreateIssuer](
					command.NewCreateIssuerHandler(repos.Issuers)),
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

				esrcwatermill.NewHandler[financing.BidOnInvoicePlacedEvent](
					eventhandler.NewBidOnInvoicePlacedHandler(repos.Invoices)),

				esrcwatermill.NewHandler[financing.BidOnInvoiceRejectedEvent](
					eventhandler.NewBidOnInvoiceRejectedHandler(repos.Investors)),

				esrcwatermill.NewHandler[financing.InvoiceFinancedEvent](
					eventhandler.NewInvoiceFinancedHandler(repos.Investors)),

				esrcwatermill.NewHandler[financing.InvoiceReversedEvent](
					eventhandler.NewInvoiceReversedHandler(repos.Investors)),

				esrcwatermill.NewHandler[financing.InvoiceApprovedEvent](
					eventhandler.NewInvoiceApprovedHandler(repos.Investors)),

				esrcwatermill.NewHandler[intevent.InvestorRegistered](
					intevent.NewInvestorRegisteredHandler(cb)),

				esrcwatermill.NewHandler[intevent.IssuerRegistered](
					intevent.NewIssuerRegisteredHandler(cb)),
			}
		},
		EventsSubscriberConstructor: func(handlerName string) (message.Subscriber, error) {
			return amqpConfig.BuildEventSubscriber(log, handlerName)
		},
		EventsPublisher:       eventsPublisher,
		Router:                router,
		CommandEventMarshaler: CommandEventMarshaler,
		Logger:                wmlog,
	})
	if err != nil {
		router.Close()
		return nil, nil, err
	}

	return facade, router, nil
}

func generateEventsTopic(eventName string) string {
	// because we are using PubSub RabbitMQ config, we can use one topic for all events
	return "events"
}

var CommandEventMarshaler = esrcwatermill.RelayEventMarshaler{
	CmdMarshaler: cqrs.JSONMarshaler{
		GenerateName: cqrs.StructName,
	},
}

func errorHandlingMiddleware(log *zap.Logger) message.HandlerMiddleware {
	return func(h message.HandlerFunc) message.HandlerFunc {
		return func(msg *message.Message) ([]*message.Message, error) {
			events, err := h(msg)
			if err != nil {
				switch err {
				case esrc.ErrOptimisticConcurrency:
					log.Warn("err handling message, retrying", zap.Error(err))
					return events, err
				default:
					log.Info("err handling message, not retrying", zap.Error(err))
					return events, nil
				}
			}

			return events, nil
		}
	}
}
