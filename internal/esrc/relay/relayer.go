package relay

import (
	"context"
	"sync/atomic"
	"time"

	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/metrics/discard"
	"go.uber.org/zap"
)

type Relayer struct {
	outbox    EventStoreOutbox
	publisher Publisher

	log                *zap.Logger
	relaysCount        metrics.Counter
	eventsRelayedCount metrics.Counter

	interval time.Duration
	stopped  int32
}

func NewRelayer(o EventStoreOutbox, p Publisher, opts ...RelayerOption) *Relayer {
	c := &Relayer{
		outbox:             o,
		publisher:          p,
		interval:           time.Millisecond * 200,
		log:                zap.NewNop(),
		relaysCount:        discard.NewCounter(),
		eventsRelayedCount: discard.NewCounter(),
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

type RelayerOption func(*Relayer)

// RelayWitInterval sets the time duration that the Relayer will wait within loops
func RelayerWithInterval(interval time.Duration) RelayerOption {
	return func(c *Relayer) {
		c.interval = interval
	}
}

func RelayerWithLogger(l *zap.Logger) RelayerOption {
	return func(c *Relayer) {
		c.log = l
	}
}

func RelayerWithMetrics(syncCallsCount, syncItemsCount metrics.Counter) RelayerOption {
	return func(c *Relayer) {
		c.relaysCount = syncCallsCount
		c.eventsRelayedCount = syncItemsCount
	}
}

func (c *Relayer) Run() {
	c.log.Info("running Relay")

	t := time.NewTicker(c.interval)
	defer t.Stop()

	for !c.isStopped() {
		<-t.C
		c.relaysCount.Add(1)

		err := c.relay(context.Background())
		if err != nil {
			c.log.Error("relay err", zap.Error(err))
		}
	}

	c.log.Info("stopping Relay")
}

func (r *Relayer) relay(ctx context.Context) error {
	events, err := r.outbox.UnpublishedEvents(ctx)
	r.log.Debug("unpublished events obtained", zap.Int("count", len(events)))
	if err != nil {
		return err
	}
	if len(events) == 0 {
		return nil
	}

	published := make([]RelayEvent, 0, len(events))
	for _, e := range events {
		err = r.publisher.Publish(ctx, e)
		if err != nil {
			r.log.Error("err publishing event", zap.Error(err), zap.Any("aggregateID", e.AggregateID), zap.Uint64("event sequence", e.Sequence))
			break
		}
		published = append(published, e)
	}
	r.log.Debug("published events", zap.Int("count", len(events)))

	if len(published) > 0 {
		err := r.outbox.MarkEventsAsPublised(ctx, events)
		if err != nil {
			return err
		}
		r.log.Debug("events have been marked as published")
	}

	return nil
}

func (c *Relayer) isStopped() bool {
	return atomic.LoadInt32(&(c.stopped)) != 0
}

func (c *Relayer) Stop() {
	atomic.StoreInt32(&(c.stopped), int32(1))
}
