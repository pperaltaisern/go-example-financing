package esrc

import (
	"context"
	"errors"
)

// EventStore persists and loads aggregates as a sequence of events
type EventStore interface {
	// Load looks for all events of an aggregate ID
	// errs:
	// 		ErrAggregateNotFound
	Load(context.Context, ID) ([]RawEvent, error)
	// Contains looks if there is an event stream for an aggregate ID without loading its events
	Contains(context.Context, ID) (bool, error)
	// Create adds a new aggregate to the event store
	// errs:
	// 		ErrAggregateAlreadyExists
	//		ErrAggregateRequiresEvents
	Create(context.Context, AggregateType, ID, []RawEvent) error
	// AppendEvents adds events to an existing event stream
	// errs:
	//		ErrAggregateNotFound
	//		ErrOptimisticConcurrency
	AppendEvents(ctx context.Context, id ID, fromVersion int, events []RawEvent) error
}

// ID is the aggregate id, modeled as empty interface since it's a specific domain concern.
// Ideally, the EventStore's implementation should be able to handle different id types by configuration.
type ID interface{}

// AggregateType is stored as part of the event stream
type AggregateType string

var (
	ErrAggregateNotFound       = errors.New("aggregate not found")
	ErrAggregateAlreadyExists  = errors.New("aggregate already exists")
	ErrAggregateRequiresEvents = errors.New("an aggregate requires events")
	ErrOptimisticConcurrency   = errors.New("optimistic concurrency error")
)
