package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/pperaltaisern/financing/pkg/config"
	"github.com/pperaltaisern/financing/pkg/financing"
	"github.com/pperaltaisern/financing/pkg/grpc"
	"github.com/pperaltaisern/financing/pkg/intevent"

	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/ThreeDotsLabs/watermill/message"
	"go.uber.org/zap"
)

func main() {
	// workarround for docker compose not waiting on dependencies correcly
	config.Wait()

	log, err := config.LoadLoggerConfig().Build()
	if err != nil {
		panic(err)
	}

	repos, _, err := config.LoadCommandPostgresConfig().BuildRepositories()
	if err != nil {
		log.Fatal("err building Postgres repositories: %v", zap.Error(err))
	}

	cqrsFacade, messageRouter, err := config.BuildCqrsFacade(log, repos)
	if err != nil {
		log.Fatal("err building CQRS facade: %v", zap.Error(err))
	}

	serverConfig := config.LoadCommandServerConfig()
	m := Main{
		log:           log,
		messageRouter: messageRouter,
		commandServer: grpc.NewCommandServer(
			serverConfig.Network,
			serverConfig.Address,
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
		PublishTestIntegrationEvents(cqrsFacade.EventBus())
	}

	log.Info("ready")
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

func PublishTestIntegrationEvents(bus *cqrs.EventBus) {
	for i := 0; i < 5; i++ {
		issuerCreated := intevent.IssuerRegistered{
			ID:   financing.NewID(),
			Name: fmt.Sprintf("ISSUER_%v", i+1),
		}
		bus.Publish(context.Background(), issuerCreated)
	}

	for i := 0; i < 5; i++ {
		investorCreated := intevent.InvestorRegistered{
			ID:      financing.NewID(),
			Name:    fmt.Sprintf("INVESTOR_%v", i+1),
			Balance: financing.Money(100),
		}
		bus.Publish(context.Background(), investorCreated)
	}
}
