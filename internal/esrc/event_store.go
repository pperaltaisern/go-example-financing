package esrc

import (
	"context"
	"errors"
)

// EventStore persists and loads aggregates as a sequence of events. Implementations must ensure optimistic concurrency.
type EventStore interface {
	// AddAggregate adds a new aggregate to the event store.
	// errs:
	// 		ErrAggregateAlreadyExists
	//		ErrAggregateRequiresEvents
	AddAggregate(context.Context, AggregateType, ID, []RawEvent) error
	// AppendEvents adds events to an existing aggregate.
	// errs:
	//		ErrAggregateNotFound
	//		ErrOptimisticConcurrency
	AppendEvents(ctx context.Context, t AggregateType, id ID, fromVersion int, events []RawEvent) error
	// AddSnapshot stores a snapshot of an aggregate.
	// errs:
	//		ErrAggregateNotFound
	//		ErrSnapshotWithGreaterVersionThanAggregate
	AddSnapshot(context.Context, AggregateType, ID, RawSnapshot) error
	// LatestSnapshot retrieves the snapshot with higher version of an aggregate. Returns a nil snapshot if there isn't any.
	LatestSnapshot(context.Context, AggregateType, ID) (*RawSnapshot, error)
	// Events retrieves all events of an aggregate from a given event number. If requested from event 0 all events
	// are obtained.
	Events(ctx context.Context, t AggregateType, id ID, fromEventNumber int) ([]RawEvent, error)
	// Contains looks if there is an event stream for an aggregate ID without loading its events,
	// returns true if the aggregate exists.
	ContainsAggregate(context.Context, AggregateType, ID) (bool, error)
}

// ID is the aggregate id, modeled as empty interface since it's a specific domain concern.
// Ideally, the EventStore's implementation should be able to handle different id types by configuration.
type ID interface{}

// AggregateType is stored as part of the event stream, needed when IDs are only unique within the same AggregateType
type AggregateType string

var (
	ErrAggregateNotFound                       = errors.New("aggregate not found")
	ErrAggregateAlreadyExists                  = errors.New("aggregate already exists")
	ErrAggregateRequiresEvents                 = errors.New("an aggregate requires events")
	ErrOptimisticConcurrency                   = errors.New("optimistic concurrency error")
	ErrSnapshotWithGreaterVersionThanAggregate = errors.New("snapshot version is greater than that of the aggregate")
)
