package main

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	kitzap "github.com/go-kit/kit/log/zap"
	"github.com/pperaltaisern/financing/internal/cmd"
	"github.com/pperaltaisern/financing/internal/esrc/esrcpg"
	"github.com/pperaltaisern/financing/internal/esrc/esrcwatermill"
	"github.com/pperaltaisern/financing/internal/esrc/relay"
	"github.com/pperaltaisern/financing/pkg/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	// workarround for docker compose not waiting on dependencies correcly
	config.Wait()

	log, err := config.LoadLoggerConfig().Build()
	if err != nil {
		panic(err)
	}

	relayCmd := relayCommand(log)
	bgRelay := cmd.NewBackground(relayCmd, cmd.BackgroundWithInterval(200*time.Millisecond))

	errC := make(chan error, 1)
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		errC <- fmt.Errorf("%s", <-c)
	}()

	go bgRelay.Run()

	log.Info("ready")
	log.Info("terminated", zap.Error(<-errC))
	bgRelay.Stop()
	log.Info("closed gracefully")
}

func relayCommand(log *zap.Logger) *relay.Command {
	pool, err := config.LoadCommandPostgresConfig().Build()
	if err != nil {
		log.Fatal("err connecting to Postgres: %v", zap.Error(err))
	}

	bus, err := config.LoadAMQPConfig().BuildEventBus(log)
	if err != nil {
		log.Fatal("err building event bus: %v", zap.Error(err))
	}

	return relay.NewCommand(
		esrcpg.NewEventStoreOutbox(pool),
		esrcwatermill.NewPublisher(bus),
		relay.CommandWithLogger(kitzap.NewZapSugarLogger(log, zapcore.DebugLevel)),
	)
}
