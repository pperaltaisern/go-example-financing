package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/pperaltaisern/financing/pkg/config"
	"github.com/pperaltaisern/financing/pkg/grpc"
	"github.com/pperaltaisern/financing/pkg/query/pg"

	"go.uber.org/zap"
)

func main() {
	// workarround for docker compose not waiting on dependencies correcly
	config.Wait()

	log, err := config.LoadLoggerConfig().Build()
	if err != nil {
		panic(err)
	}

	db, err := config.LoadQueryPostgresConfig().BuildGORM()
	if err != nil {
		log.Fatal("error connecting to postgres", zap.Error(err))
	}

	serverConfig := config.LoadQueryServerConfig()
	m := Main{
		log: log,
		queryServer: grpc.NewQueryServer(
			serverConfig.Network,
			serverConfig.Address,
			pg.NewInvestorQueries(db),
			pg.NewInvoiceQueries(db),
			pg.NewIssuerQueries(db),
		),
	}

	errC := make(chan error, 4)
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		errC <- fmt.Errorf("%s", <-c)
	}()

	m.Run(errC)

	log.Info("ready")
	log.Info("terminated", zap.Error(<-errC))
	m.Close()
	log.Info("closed gracefully")
}

type Main struct {
	log         *zap.Logger
	queryServer *grpc.QueryServer
}

func (m *Main) Run(errC chan<- error) {
	go func() { errC <- m.queryServer.Open() }()
}

func (m *Main) Close() {
	m.queryServer.Close()
	m.log.Sync()
}
