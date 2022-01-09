package esrc

import (
	"context"
)

// Repository is a layer on top of an EventStore that reduces the boilerplate needed to build a domain repository
type Repository struct {
	aggregateType  AggregateType
	eventStore     EventStore
	eventFactory   EventFactory
	eventMarshaler EventMarshaler
}

func NewRepository(at AggregateType, es EventStore, ef EventFactory, em EventMarshaler) Repository {
	return Repository{
		aggregateType:  at,
		eventStore:     es,
		eventFactory:   ef,
		eventMarshaler: em,
	}
}

func (r Repository) ByID(ctx context.Context, id ID) ([]Event, error) {
	rawEvents, err := r.eventStore.Load(ctx, id)
	if err != nil {
		return nil, err
	}

	events := make([]Event, len(rawEvents))
	for i, raw := range rawEvents {
		e, err := r.eventFactory.CreateEmptyEvent(raw.Name)
		if err != nil {
			return nil, err
		}
		err = r.eventMarshaler.UnmarshalEvent(raw.Data, e)
		if err != nil {
			return nil, err
		}
		events[i] = e
	}

	return events, nil
}

func (r Repository) Contains(ctx context.Context, id ID) (bool, error) {
	return r.eventStore.Contains(ctx, id)
}

func (r Repository) Update(ctx context.Context, id ID, fromVersion int, events []Event) error {
	if len(events) == 0 {
		return nil
	}

	rawEvents, err := MarshalEvents(events, r.eventMarshaler)
	if err != nil {
		return err
	}
	return r.eventStore.AppendEvents(ctx, id, fromVersion, rawEvents)
}

func (r Repository) Add(ctx context.Context, id ID, events []Event) error {
	rawEvents, err := MarshalEvents(events, r.eventMarshaler)
	if err != nil {
		return err
	}
	return r.eventStore.Create(ctx, r.aggregateType, id, rawEvents)
}
