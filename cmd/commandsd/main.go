package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/pperaltaisern/financing/pkg/command"
	"github.com/pperaltaisern/financing/pkg/config"
	"github.com/pperaltaisern/financing/pkg/financing"
	"github.com/pperaltaisern/financing/pkg/grpc"

	"github.com/ThreeDotsLabs/watermill/message"
	"go.uber.org/zap"
)

func main() {
	log, err := config.LoadLoggerConfig().Build()
	if err != nil {
		panic(err)
	}

	repos, err := config.LoadPostgresConfig().BuildRepositories()
	if err != nil {
		log.Fatal("err building Postgres repositories: %v", zap.Error(err))
	}

	cqrsFacade, messageRouter, err := config.BuildCqrsFacade(log, repos)
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

type Main struct {
	log           *zap.Logger
	messageRouter *message.Router
	commandServer *grpc.CommandServer
}

func (m *Main) Run(errC chan<- error) {
	go func() { errC <- m.messageRouter.Run(context.Background()) }()
	go func() { errC <- m.commandServer.Open() }()
}

func (m *Main) Close() {
	err := m.messageRouter.Close()
	if err != nil {
		m.log.Error("err clossing message router %v: err")
	}
	m.commandServer.Close()
	m.log.Sync()
}
