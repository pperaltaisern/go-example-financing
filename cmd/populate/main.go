package main

import (
	"context"
	"fmt"
	"log"

	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/pperaltaisern/financing/pkg/config"
	"github.com/pperaltaisern/financing/pkg/financing"
	"github.com/pperaltaisern/financing/pkg/intevent"
	"go.uber.org/zap"
)

func main() {
	logger, err := config.LoadLoggerConfig().Build()
	if err != nil {
		log.Fatalf("err instancing logger: %v", err)
	}

	bus, err := config.LoadAMQPConfig().BuildEventBus(logger)
	if err != nil {
		logger.Fatal("err building bus", zap.Error(err))
	}
	PublishTestIntegrationEvents(bus)
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
