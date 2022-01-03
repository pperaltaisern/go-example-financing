package esrc

import "context"

type EventStore interface {
	Load(context.Context, ID) ([]Event, error)
	Contains(context.Context, ID) (bool, error)
	Create(context.Context, Aggregate, AggregateType) error
	AppendEvents(context.Context, Aggregate) error
}

// AggregateType is saved as part of the event stream in the event store
type AggregateType string
