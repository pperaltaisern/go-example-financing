package esrctesting

import (
	"context"
	"ledger/internal/esrc"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type EventStoreAcceptance struct {
	es    esrc.EventStore
	newID func() esrc.ID
}

func NewEventStoreAcceptance(es esrc.EventStore, opts ...EventStoreAcceptanceOption) *EventStoreAcceptance {
	a := &EventStoreAcceptance{
		es:    es,
		newID: func() esrc.ID { return uuid.New() },
	}

	for _, opt := range opts {
		opt(a)
	}
	return a
}

type EventStoreAcceptanceOption func(*EventStoreAcceptance)

func EventStoreAcceptanceNewID(newID func() esrc.ID) EventStoreAcceptanceOption {
	return func(a *EventStoreAcceptance) {
		a.newID = newID
	}
}

func (a *EventStoreAcceptance) Test(t *testing.T) {
	t.Parallel()
	t.Run("TestFromEmptyEventStore", func(t *testing.T) {
		t.Parallel()
		t.Run("LoadNotExistingAggregate", func(t *testing.T) {
			t.Parallel()
			a.AssertLoadNotExistingAggregate(t)
		})
		t.Run("NotContains", func(t *testing.T) {
			t.Parallel()
			a.AssertContainsNotExistingAggregate(t)
		})
		t.Run("CreateEmptyAggregate", func(t *testing.T) {
			t.Parallel()
			a.AssertCreateEmptyAggregate(t)
		})
		t.Run("AppendEventsToNotExistingAggregate", func(t *testing.T) {
			t.Parallel()
			a.AssertAppendEventsToNotExistingAggregate(t)
		})
	})

	t.Run("TestFromPopulatedStore", func(t *testing.T) {
		t.Parallel()

		id := a.newID()
		initialEvents := []esrc.RawEvent{
			{
				Name: "event 1",
				Data: []byte("data 1"),
			}, {
				Name: "event 2",
				Data: []byte("data 2"),
			},
		}
		t.Run("CreateValidAggregate", func(t *testing.T) {
			err := a.es.Create(context.Background(), "type", id, initialEvents)
			assert.NoError(t, err)
		})
		t.Run("CreateAlreadyExistingAggregate", func(t *testing.T) {
			err := a.es.Create(context.Background(), "type", id, initialEvents)
			_ = assert.Error(t, err) &&
				assert.Equal(t, esrc.ErrAggregateAlreadyExists, err)
		})
		t.Run("Contains", func(t *testing.T) {
			a.AssertContaintsExistingAggregate(t, id)
		})
		t.Run("Load", func(t *testing.T) {
			loadedEvents, err := a.es.Load(context.Background(), id)
			_ = assert.NoError(t, err) &&
				assert.Equal(t, initialEvents, loadedEvents)
		})

		appendedEvents := []esrc.RawEvent{
			{
				Name: "event 3",
				Data: []byte("data 3"),
			}, {
				Name: "event 4",
				Data: []byte("data 4"),
			},
		}
		t.Run("AppendEvents", func(t *testing.T) {
			err := a.es.AppendEvents(context.Background(), id, len(initialEvents), appendedEvents)
			_ = assert.NoError(t, err)
		})
		t.Run("AppendEvents same events (simulation for optimistic concurrency)", func(t *testing.T) {
			err := a.es.AppendEvents(context.Background(), id, len(initialEvents), appendedEvents)
			_ = assert.Error(t, err) &&
				assert.Equal(t, esrc.ErrOptimisticConcurrency, err)
		})
		t.Run("Load after appended events", func(t *testing.T) {
			loadedEvents, err := a.es.Load(context.Background(), id)
			_ = assert.NoError(t, err) &&
				assert.Len(t, loadedEvents, 4) &&
				assert.Equal(t, initialEvents, loadedEvents[0:2]) &&
				assert.Equal(t, appendedEvents, loadedEvents[2:4])
		})
		t.Run("LoadNotExistingAggregate", func(t *testing.T) {
			t.Parallel()
			a.AssertLoadNotExistingAggregate(t)
		})
		t.Run("NotContains", func(t *testing.T) {
			t.Parallel()
			a.AssertContainsNotExistingAggregate(t)
		})
	})
}

func (a *EventStoreAcceptance) AssertLoadNotExistingAggregate(t *testing.T) bool {
	events, err := a.es.Load(context.Background(), a.newID())
	return assert.Error(t, err) &&
		assert.Equal(t, err, esrc.ErrAggregateNotFound) &&
		assert.Empty(t, events)
}

func (a *EventStoreAcceptance) AssertContaintsExistingAggregate(t *testing.T, id esrc.ID) bool {
	found, err := a.es.Contains(context.Background(), id)
	return assert.NoError(t, err) &&
		assert.True(t, found)
}

func (a *EventStoreAcceptance) AssertContainsNotExistingAggregate(t *testing.T) bool {
	found, err := a.es.Contains(context.Background(), a.newID())
	return assert.NoError(t, err) &&
		assert.False(t, found)
}

func (a *EventStoreAcceptance) AssertCreateEmptyAggregate(t *testing.T) bool {
	events := make([]esrc.RawEvent, 0)

	err := a.es.Create(context.Background(), "type", a.newID(), events)
	return assert.Error(t, err) &&
		assert.Equal(t, err, esrc.ErrAggregateRequiresEvents)
}

func (a *EventStoreAcceptance) AssertAppendEventsToNotExistingAggregate(t *testing.T) bool {
	events := make([]esrc.RawEvent, 1)

	err := a.es.AppendEvents(context.Background(), a.newID(), 0, events)
	return assert.Error(t, err) &&
		assert.Equal(t, err, esrc.ErrAggregateNotFound)
}
