package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync/atomic"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/pperaltaisern/financing/pkg/config"
	"github.com/pperaltaisern/financing/pkg/projection"
	"github.com/pperaltaisern/financing/pkg/query/pg"
	"go.uber.org/zap"
)

func main() {
	config.Wait()

	log, err := config.LoadLoggerConfig().Build()
	if err != nil {
		panic(err)
	}

	subscriber, err := config.LoadAMQPConfig().BuildEventSubscriber(log, "query")
	if err != nil {
		log.Fatal("err building event bus: %v", zap.Error(err))
	}

	db, err := config.LoadQueryPostgresConfig().BuildGORM()
	if err != nil {
		log.Fatal("err connecting to Postgres: %v", zap.Error(err))
	}

	eventProjector, err := pg.NewEventProjector(db)
	if err != nil {
		log.Fatal("err building Postgres event projector: %v", zap.Error(err))
	}
	messageProjector := projection.NewMessageProjector(
		eventProjector,
		config.CommandEventMarshaler,
		func(m *message.Message, e error) {
			log.Error("projection err", zap.Error(e), zap.Any("message", m))
		})

	eventHandler := NewProjectorSubscriber(subscriber, messageProjector)

	errC := make(chan error, 2)
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		errC <- fmt.Errorf("%s", <-c)
	}()

	go func() { errC <- eventHandler.Run() }()

	log.Info("ready")
	log.Info("terminated", zap.Error(<-errC))
	eventHandler.Stop()
	log.Info("closed gracefully")
}

type ProjectorSubscriber struct {
	subscriber       message.Subscriber
	messageProjector *projection.MessageProjector
	close            int32
}

func NewProjectorSubscriber(sub message.Subscriber, mp *projection.MessageProjector) *ProjectorSubscriber {
	return &ProjectorSubscriber{
		subscriber:       sub,
		messageProjector: mp,
	}
}

func (h *ProjectorSubscriber) Run() error {
	messageC, err := h.subscriber.Subscribe(context.Background(), "events")
	if err != nil {
		return err
	}
	for !h.isClosed() {
		m := <-messageC
		h.messageProjector.ProjectMessage(m)
	}
	return nil
}

func (h *ProjectorSubscriber) isClosed() bool {
	return h.close > 0
}

func (h *ProjectorSubscriber) Stop() {
	h.subscriber.Close()
	atomic.AddInt32(&h.close, 1)
}
