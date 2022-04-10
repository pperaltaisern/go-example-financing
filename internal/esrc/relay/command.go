package relay

import (
	"context"

	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/metrics/discard"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/pperaltaisern/financing/internal/cmd"
)

type Command struct {
	outbox             EventStoreOutbox
	publisher          Publisher
	log                log.Logger
	eventsRelayedCount metrics.Counter
}

var _ cmd.Command = (*Command)(nil)

func NewCommand(o EventStoreOutbox, p Publisher, opts ...CommandOption) *Command {
	c := &Command{
		outbox:             o,
		publisher:          p,
		log:                log.NewNopLogger(),
		eventsRelayedCount: discard.NewCounter(),
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

type CommandOption func(*Command)

func CommandWithLogger(l log.Logger) CommandOption {
	return func(c *Command) {
		c.log = l
	}
}

func CommandWithMetrics(eventsRelayedCount metrics.Counter) CommandOption {
	return func(c *Command) {
		c.eventsRelayedCount = eventsRelayedCount
	}
}

func (r *Command) Execute(ctx context.Context) error {
	events, err := r.outbox.UnpublishedEvents(ctx)
	if err != nil {
		level.Error(r.log).Log("msg", "unpublished events error", "err", err)
		return err
	}
	if len(events) == 0 {
		return nil
	}
	level.Info(r.log).Log("msg", "unpublished events obtained", "count", len(events))

	published := make([]RelayEvent, 0, len(events))
	for _, e := range events {
		err = r.publisher.Publish(ctx, e)
		if err != nil {
			level.Error(r.log).Log(
				"msg", "publish events error",
				"err", err,
				"aggregateID", e.AggregateID,
				"event sequence", e.Sequence)
			break
		}
		published = append(published, e)
	}
	level.Info(r.log).Log("msg", "published events", "count", len(events))
	if len(published) == 0 {
		return nil
	}

	err = r.outbox.MarkEventsAsPublised(ctx, events)
	if err != nil {
		level.Error(r.log).Log("msg", "mark events as published error", "err", err)
		return err
	}
	level.Info(r.log).Log("msg", "events marked as published", "count", len(events))
	r.eventsRelayedCount.Add(float64(len(events)))
	return nil
}
