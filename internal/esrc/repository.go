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

func NewRepository(t AggregateType, es EventStore, ef EventFactory, em EventMarshaler) *Repository {
	return &Repository{
		aggregateType:  t,
		eventStore:     es,
		eventFactory:   ef,
		eventMarshaler: em,
	}
}

func (r *Repository) FindByID(ctx context.Context, id ID) ([]Event, error) {
	rawEvents, err := r.eventStore.Load(ctx, r.aggregateType, id)
	if err != nil {
		return nil, err
	}

	events := make([]Event, len(rawEvents))
	for i, raw := range rawEvents {
		event, err := r.eventFactory.CreateEmptyEvent(raw.Name)
		if err != nil {
			return nil, err
		}
		err = r.eventMarshaler.UnmarshalEvent(raw.Data, event)
		if err != nil {
			return nil, err
		}
		events[i] = event
	}

	return events, nil
}

func (r *Repository) Contains(ctx context.Context, id ID) (bool, error) {
	return r.eventStore.Contains(ctx, r.aggregateType, id)
}

func (r *Repository) Update(ctx context.Context, id ID, fromVersion int, events []Event) error {
	if len(events) == 0 {
		return nil
	}

	rawEvents, err := MarshalEvents(events, r.eventMarshaler)
	if err != nil {
		return err
	}
	return r.eventStore.AppendEvents(ctx, r.aggregateType, id, fromVersion, rawEvents)
}

func (r *Repository) Add(ctx context.Context, id ID, events []Event) error {
	rawEvents, err := MarshalEvents(events, r.eventMarshaler)
	if err != nil {
		return err
	}
	return r.eventStore.Create(ctx, r.aggregateType, id, rawEvents)
}
