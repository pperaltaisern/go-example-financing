package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/pperaltaisern/financing/internal/esrc/esrcpg"
	"github.com/pperaltaisern/financing/internal/esrc/esrcwatermill"
	"github.com/pperaltaisern/financing/internal/esrc/relay"
	"github.com/pperaltaisern/financing/pkg/config"
	"go.uber.org/zap"
)

func main() {
	// workarround for docker compose not waiting on dependencies correcly
	config.Wait()

	log, err := config.LoadLoggerConfig().Build()
	if err != nil {
		panic(err)
	}

	pool, err := config.LoadPostgresConfig().Build()
	if err != nil {
		log.Fatal("err connecting to Postgres: %v", zap.Error(err))
	}

	bus, err := config.BuildEventBus(log)
	if err != nil {
		log.Fatal("err building event bus: %v", zap.Error(err))
	}

	relayer := relay.NewRelayer(
		esrcpg.NewEventStoreOutbox(pool),
		esrcwatermill.NewPublisher(bus),
		relay.RelayerWithLogger(log),
	)

	errC := make(chan error, 1)
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		errC <- fmt.Errorf("%s", <-c)
	}()

	go relayer.Run()

	log.Info("ready")
	log.Info("terminated", zap.Error(<-errC))
	relayer.Stop()
	log.Info("closed gracefully")
}
