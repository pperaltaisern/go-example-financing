package esrctesting

import (
	"context"
	"testing"

	"github.com/pperaltaisern/financing/internal/esrc"
	"github.com/pperaltaisern/financing/internal/esrc/relay"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type EventStoreAcceptanceSuite struct {
	eventStore esrc.EventStore
	outbox     relay.EventStoreOutbox
	newID      func() esrc.ID
}

func NewEventStoreAcceptanceSuite(es esrc.EventStore, opts ...EventStoreAcceptanceSuiteOption) *EventStoreAcceptanceSuite {
	a := &EventStoreAcceptanceSuite{
		eventStore: es,
		newID:      func() esrc.ID { return uuid.New() },
	}

	for _, opt := range opts {
		opt(a)
	}
	return a
}

type EventStoreAcceptanceSuiteOption func(*EventStoreAcceptanceSuite)

func EventStoreAcceptanceSuiteNewID(newID func() esrc.ID) EventStoreAcceptanceSuiteOption {
	return func(a *EventStoreAcceptanceSuite) {
		a.newID = newID
	}
}

func EventStoreAcceptanceSuiteWithOutbox(o relay.EventStoreOutbox) EventStoreAcceptanceSuiteOption {
	return func(a *EventStoreAcceptanceSuite) {
		a.outbox = o
	}
}

const testAggregateType = "test"

func (a *EventStoreAcceptanceSuite) Test(t *testing.T) {
	t.Parallel()
	t.Run("TestFromEmptyEventStore", func(t *testing.T) {
		t.Run("LoadNotExistingAggregate", func(t *testing.T) {
			a.AssertLoadNotExistingAggregate(t)
		})
		t.Run("NotContains", func(t *testing.T) {
			a.AssertContainsNotExistingAggregate(t)
		})
		t.Run("CreateEmptyAggregate", func(t *testing.T) {
			a.AssertCreateEmptyAggregate(t)
		})
		t.Run("AppendEventsToNotExistingAggregate", func(t *testing.T) {
			a.AssertAppendEventsToNotExistingAggregate(t)
		})
	})

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
	appendedEvents := []esrc.RawEvent{
		{
			Name: "event 3",
			Data: []byte("data 3"),
		}, {
			Name: "event 4",
			Data: []byte("data 4"),
		},
	}

	t.Run("TestFromPopulatedStore", func(t *testing.T) {
		t.Run("CreateValidAggregate", func(t *testing.T) {
			err := a.eventStore.AddAggregate(context.Background(), "type", id, initialEvents)
			assert.NoError(t, err)
		})
		t.Run("CreateAlreadyExistingAggregate", func(t *testing.T) {
			err := a.eventStore.AddAggregate(context.Background(), "type", id, initialEvents)
			_ = assert.Error(t, err) &&
				assert.Equal(t, esrc.ErrAggregateAlreadyExists, err)
		})
		t.Run("Contains", func(t *testing.T) {
			a.AssertContaintsExistingAggregate(t, id)
		})
		t.Run("Load", func(t *testing.T) {
			loadedEvents, err := a.eventStore.Events(context.Background(), testAggregateType, id, 0)
			_ = assert.NoError(t, err) &&
				assert.Equal(t, initialEvents, loadedEvents)
		})

		t.Run("AppendEvents", func(t *testing.T) {
			err := a.eventStore.AppendEvents(context.Background(), testAggregateType, id, len(initialEvents), appendedEvents)
			_ = assert.NoError(t, err)
		})
		t.Run("AppendEvents same events (simulation for optimistic concurrency)", func(t *testing.T) {
			err := a.eventStore.AppendEvents(context.Background(), testAggregateType, id, len(initialEvents), appendedEvents)
			_ = assert.Error(t, err) &&
				assert.Equal(t, esrc.ErrOptimisticConcurrency, err)
		})
		t.Run("Load after appended events", func(t *testing.T) {
			loadedEvents, err := a.eventStore.Events(context.Background(), testAggregateType, id, 0)
			_ = assert.NoError(t, err) &&
				assert.Len(t, loadedEvents, 4) &&
				assert.Equal(t, initialEvents, loadedEvents[0:2]) &&
				assert.Equal(t, appendedEvents, loadedEvents[2:4])
		})
		t.Run("LoadNotExistingAggregate", func(t *testing.T) {
			a.AssertLoadNotExistingAggregate(t)
		})
		t.Run("NotContains", func(t *testing.T) {
			a.AssertContainsNotExistingAggregate(t)
		})
	})

	if a.outbox != nil {
		t.Run("TestEventStoreOutbox", func(t *testing.T) {
			expectedUnpublishedEvents := make([]relay.RelayEvent, 4)
			for i := 0; i < 4; i++ {
				var e esrc.RawEvent
				if i < 2 {
					e = initialEvents[i]
				} else {
					e = appendedEvents[i-2]
				}
				expectedUnpublishedEvents[i] = relay.NewRelayEvent(id, uint64(i+1), e)
			}

			t.Run("UnpublishedEvents", func(t *testing.T) {
				events, err := a.outbox.UnpublishedEvents(context.Background())
				_ = assert.NoError(t, err) &&
					assert.Len(t, events, 4) &&
					assert.Equal(t, expectedUnpublishedEvents, events)
			})
			t.Run("MarkAsPublished", func(t *testing.T) {
				err := a.outbox.MarkEventsAsPublised(context.Background(), expectedUnpublishedEvents)
				_ = assert.NoError(t, err)
			})
			t.Run("UnpublishedEvents after publishing", func(t *testing.T) {
				events, err := a.outbox.UnpublishedEvents(context.Background())
				_ = assert.NoError(t, err) &&
					assert.Empty(t, events)
			})
		})
	}
}

func (a *EventStoreAcceptanceSuite) AssertLoadNotExistingAggregate(t *testing.T) bool {
	events, err := a.eventStore.Events(context.Background(), testAggregateType, a.newID(), 0)
	return assert.Error(t, err) &&
		assert.Equal(t, err, esrc.ErrAggregateNotFound) &&
		assert.Empty(t, events)
}

func (a *EventStoreAcceptanceSuite) AssertContaintsExistingAggregate(t *testing.T, id esrc.ID) bool {
	found, err := a.eventStore.ContainsAggregate(context.Background(), testAggregateType, id)
	return assert.NoError(t, err) &&
		assert.True(t, found)
}

func (a *EventStoreAcceptanceSuite) AssertContainsNotExistingAggregate(t *testing.T) bool {
	found, err := a.eventStore.ContainsAggregate(context.Background(), testAggregateType, a.newID())
	return assert.NoError(t, err) &&
		assert.False(t, found)
}

func (a *EventStoreAcceptanceSuite) AssertCreateEmptyAggregate(t *testing.T) bool {
	events := make([]esrc.RawEvent, 0)

	err := a.eventStore.AddAggregate(context.Background(), "type", a.newID(), events)
	return assert.Error(t, err) &&
		assert.Equal(t, err, esrc.ErrAggregateRequiresEvents)
}

func (a *EventStoreAcceptanceSuite) AssertAppendEventsToNotExistingAggregate(t *testing.T) bool {
	events := make([]esrc.RawEvent, 1)

	err := a.eventStore.AppendEvents(context.Background(), testAggregateType, a.newID(), 0, events)
	return assert.Error(t, err) &&
		assert.Equal(t, err, esrc.ErrAggregateNotFound)
}
