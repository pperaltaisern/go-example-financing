package esrc

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

// Repository is a layer on top of an EventStore that reduces the boilerplate needed to build an event sourced aggregate repository
type Repository[T Aggregate] struct {
	eventStore        EventStore
	aggregateType     AggregateType
	aggregateFactory  AggregateFactory[T]
	eventFactory      EventFactory
	eventMarshaler    EventMarshaler
	eventsPerSnapshot int
	log               log.Logger
}

func NewRepository[T Aggregate](es EventStore, af AggregateFactory[T], ef EventFactory, opts ...RepositoryOption[T]) *Repository[T] {
	r := &Repository[T]{
		eventStore:        es,
		aggregateType:     aggregateTypeFromAggregate[T](),
		aggregateFactory:  af,
		eventFactory:      ef,
		eventMarshaler:    JSONEventMarshaler{},
		eventsPerSnapshot: 200,
		log:               log.NewNopLogger(),
	}
	for _, opt := range opts {
		opt(r)
	}

	return r
}

func aggregateTypeFromAggregate[T Aggregate]() AggregateType {
	var aggregate T
	t := fmt.Sprintf("%T", aggregate)
	t = strings.Split(t, ".")[1]
	return AggregateType(t)
}

type RepositoryOption[T Aggregate] func(*Repository[T])

func RepositoryWithEventsMarshaler[T Aggregate](eventMarshaler EventMarshaler) RepositoryOption[T] {
	return func(r *Repository[T]) {
		r.eventMarshaler = eventMarshaler
	}
}

func RepositoryWithEventsPerSnapshot[T Aggregate](eventsPerSnapshot int) RepositoryOption[T] {
	return func(r *Repository[T]) {
		r.eventsPerSnapshot = eventsPerSnapshot
	}
}

func RepositoryWithLogger[T Aggregate](log log.Logger) RepositoryOption[T] {
	return func(r *Repository[T]) {
		r.log = log
	}
}

func (r *Repository[T]) FindByID(ctx context.Context, id ID) (T, error) {
	var nilAggregate T

	snapshot, err := r.eventStore.LatestSnapshot(ctx, r.aggregateType, id)
	if err != nil {
		return nilAggregate, err
	}
	fromEventVersion := 0
	if snapshot != nil {
		fromEventVersion = snapshot.Version + 1
	}

	rawEvents, err := r.eventStore.Events(ctx, r.aggregateType, id, fromEventVersion)
	if err != nil {
		return nilAggregate, err
	}

	events := make([]Event, len(rawEvents))
	for i, raw := range rawEvents {
		event, err := r.eventFactory.CreateEmptyEvent(raw.Name)
		if err != nil {
			return nilAggregate, err
		}
		err = r.eventMarshaler.UnmarshalEvent(raw.Data, event)
		if err != nil {
			return nilAggregate, err
		}
		events[i] = event
	}

	if snapshot != nil {
		return r.aggregateFactory.NewAggregateFromSnapshotAndEvents(*snapshot, events)
	}
	return r.aggregateFactory.NewAggregateFromEvents(events)
}

func (r *Repository[T]) Contains(ctx context.Context, id ID) (bool, error) {
	return r.eventStore.ContainsAggregate(ctx, r.aggregateType, id)
}

func (r *Repository[T]) Update(ctx context.Context, aggregate T) error {
	events := aggregate.Changes()
	if len(events) == 0 {
		return nil
	}

	rawEvents, err := MarshalEvents(events, r.eventMarshaler)
	if err != nil {
		return err
	}
	err = r.eventStore.AppendEvents(ctx, r.aggregateType, aggregate.ID(), aggregate.InitialVersion(), rawEvents)
	if err != nil {
		return err
	}
	r.addSnapshotIfRequired(ctx, aggregate)
	return nil
}

func (r *Repository[T]) Add(ctx context.Context, aggregate T) error {
	rawEvents, err := MarshalEvents(aggregate.Changes(), r.eventMarshaler)
	if err != nil {
		return err
	}
	err = r.eventStore.AddAggregate(ctx, r.aggregateType, aggregate.ID(), rawEvents)
	if err != nil {
		return err
	}
	r.addSnapshotIfRequired(ctx, aggregate)
	return nil
}

func (r *Repository[T]) addSnapshotIfRequired(ctx context.Context, aggregate T) {
	if !r.shouldDoSnapshot(aggregate) {
		return
	}

	eventsLen := len(aggregate.Changes())
	initialVersion := aggregate.InitialVersion()
	snapshot, err := aggregate.Snapshot()
	if err != nil {
		level.Error(r.log).Log("msg", "create snapshot for aggregate error", "err", err)
		return
	}
	go func(id ID) {
		rawSnapshot := RawSnapshot{
			Version: initialVersion + eventsLen,
			Data:    snapshot,
		}
		err := r.eventStore.AddSnapshot(context.Background(), r.aggregateType, id, rawSnapshot)
		if err != nil {
			level.Error(r.log).Log("msg", "add snapshot to event store error", "err", err)
		}
	}(aggregate.ID())
}

func (r *Repository[T]) shouldDoSnapshot(aggregate T) bool {
	if r.eventsPerSnapshot <= 0 {
		return false
	}
	return r.eventsPerSnapshot > 0 &&
		(aggregate.InitialVersion()%r.eventsPerSnapshot)+len(aggregate.Changes()) >= r.eventsPerSnapshot
}
