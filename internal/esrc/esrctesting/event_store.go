package esrctesting

import (
	"context"
	"testing"

	"github.com/pperaltaisern/financing/internal/esrc"
	"github.com/pperaltaisern/financing/internal/esrc/relay"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	t.Run("WhenEventStoreIsEmpty", func(t *testing.T) {
		t.Run("ReadingEventsShouldNotFindTheAggregate", func(t *testing.T) {
			a.AssertLoadNotExistingAggregate(t)
		})
		t.Run("ContainsShouldNotFindTheAggregate", func(t *testing.T) {
			a.AssertContainsNotExistingAggregate(t)
		})
		t.Run("CreateAggregateWithoutEventsShouldReturnError", func(t *testing.T) {
			a.AssertCreateEmptyAggregate(t)
		})
		t.Run("AppendEventsToNotExistingAggregateShouldReturnError", func(t *testing.T) {
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

	t.Run("WhenPopulatingWithAggregatesAndEvents", func(t *testing.T) {
		t.Run("CreatingAValidAggregateShouldSucceed", func(t *testing.T) {
			err := a.eventStore.AddAggregate(context.Background(), "type", id, initialEvents)
			assert.NoError(t, err)
		})
		t.Run("CreatingAnAlreadyExistingAggregateShouldReturnExistingAggregateError", func(t *testing.T) {
			err := a.eventStore.AddAggregate(context.Background(), "type", id, initialEvents)
			_ = assert.Error(t, err) &&
				assert.Equal(t, esrc.ErrAggregateAlreadyExists, err)
		})
		t.Run("ContainsShouldReturnTrueForExistingAggregate", func(t *testing.T) {
			a.AssertContaintsExistingAggregate(t, id)
		})
		t.Run("ReadEventsShouldReturnAddedEventsOfExistingAggregate", func(t *testing.T) {
			loadedEvents, err := a.eventStore.Events(context.Background(), testAggregateType, id, 0)
			_ = assert.NoError(t, err) &&
				assert.Equal(t, initialEvents, loadedEvents)
		})

		t.Run("AppendEventsToExistingAggregateShouldSucceed", func(t *testing.T) {
			err := a.eventStore.AppendEvents(context.Background(), testAggregateType, id, len(initialEvents), appendedEvents)
			_ = assert.NoError(t, err)
		})
		t.Run("AppendSameEventsShouldReturnOptimisticConcurrencyErr", func(t *testing.T) {
			err := a.eventStore.AppendEvents(context.Background(), testAggregateType, id, len(initialEvents), appendedEvents)
			_ = assert.Error(t, err) &&
				assert.Equal(t, esrc.ErrOptimisticConcurrency, err)
		})
		t.Run("LoadEventsForAnAggregateWithEventsShouldReturnThem", func(t *testing.T) {
			loadedEvents, err := a.eventStore.Events(context.Background(), testAggregateType, id, 0)
			_ = assert.NoError(t, err) &&
				assert.Len(t, loadedEvents, 4) &&
				assert.Equal(t, initialEvents, loadedEvents[0:2]) &&
				assert.Equal(t, appendedEvents, loadedEvents[2:4])
		})
		t.Run("LoadNotExistingAggregateShouldNotFindTheAggregate", func(t *testing.T) {
			a.AssertLoadNotExistingAggregate(t)
		})
		t.Run("ContainsShouldNotFindTheAggregate", func(t *testing.T) {
			a.AssertContainsNotExistingAggregate(t)
		})
	})

	t.Run("WhenPopulatingWithSnapshots", func(t *testing.T) {
		snapshotVersion := 2
		snapshotVersion2 := esrc.RawSnapshot{
			Version: snapshotVersion,
			Data:    []byte("snapshot data"),
		}
		idWithSnapshot := a.newID()
		t.Run("AddSnapshotShouldReturnErrorIfAggregateIsntStored", func(t *testing.T) {
			err := a.eventStore.AddSnapshot(context.Background(), testAggregateType, a.newID(), snapshotVersion2)
			require.Equal(t, esrc.ErrAggregateNotFound, err)
		})
		t.Run("AddSnapshotWithVersion2ToExistingAggregateWith2EventsShouldSuccess", func(t *testing.T) {
			err := a.eventStore.AddAggregate(context.Background(), testAggregateType, idWithSnapshot, initialEvents)
			require.Nil(t, err)

			err = a.eventStore.AddSnapshot(context.Background(), testAggregateType, a.newID(), snapshotVersion2)
			require.Equal(t, esrc.ErrAggregateNotFound, err)
		})
		t.Run("AddSnapshotWithVersion3ToAggregateWith2EventsShouldReturnErr", func(t *testing.T) {
			id := a.newID()
			err := a.eventStore.AddAggregate(context.Background(), testAggregateType, id, initialEvents)
			require.NoError(t, err)

			invalidSnapshot := esrc.RawSnapshot{
				Version: 3,
				Data:    snapshotVersion2.Data,
			}
			err = a.eventStore.AddSnapshot(context.Background(), testAggregateType, id, invalidSnapshot)
			require.Equal(t, esrc.ErrSnapshotWithGreaterVersionThanAggregate, err)
		})
		t.Run("LatestSnapshotForANonExistingAggregateShouldReturnNilAndNoError", func(t *testing.T) {
			snapshot, err := a.eventStore.LatestSnapshot(context.Background(), testAggregateType, a.newID())
			require.Nil(t, snapshot)
			require.NoError(t, err)
		})

		t.Run("LatestSnapshotForAnExistingAggregateWithoutSnapshotShouldReturnNilAndNoError", func(t *testing.T) {
			id := a.newID()
			err := a.eventStore.AddAggregate(context.Background(), testAggregateType, id, initialEvents)
			require.NoError(t, err)

			snapshot, err := a.eventStore.LatestSnapshot(context.Background(), testAggregateType, a.newID())
			require.Nil(t, snapshot)
			require.NoError(t, err)
		})

		t.Run("LatestSnapshotForAnExistingAggregateWithMultipleSnapshotsShouldReturnTheLatest", func(t *testing.T) {
			snapshotVersion1 := esrc.RawSnapshot{
				Version: 1,
				Data:    []byte("version 1 snapshot"),
			}
			err := a.eventStore.AddSnapshot(context.Background(), testAggregateType, idWithSnapshot, snapshotVersion1)
			require.NoError(t, err)

			err = a.eventStore.AddSnapshot(context.Background(), testAggregateType, idWithSnapshot, snapshotVersion2)
			require.NoError(t, err)

			snapshot, err := a.eventStore.LatestSnapshot(context.Background(), testAggregateType, idWithSnapshot)
			require.NoError(t, err)
			require.NotNil(t, snapshot)
			require.Equal(t, snapshotVersion2.Version, snapshot.Version)
			require.Equal(t, string(snapshotVersion2.Data), string(snapshot.Data))
		})
	})

	if a.outbox != nil {
		t.Run("WhenRelayingEventsInTheOutbox", func(t *testing.T) {
			var eventsFromPreviousTests []relay.RelayEvent
			t.Run("ReadingExistingUnpublishedEventsShouldReturnEventsFromPreviousTests(cleaning)", func(t *testing.T) {
				var err error
				eventsFromPreviousTests, err = a.outbox.UnpublishedEvents(context.Background())
				require.NoError(t, err)
				require.NotEmpty(t, eventsFromPreviousTests)

				err = a.outbox.MarkEventsAsPublised(context.Background(), eventsFromPreviousTests)
				require.NoError(t, err)
			})
			t.Run("MarkingEventsFromPreviousTestsAsPublishedShouldSucceed(cleaning)", func(t *testing.T) {
				err := a.outbox.MarkEventsAsPublised(context.Background(), eventsFromPreviousTests)
				require.NoError(t, err)
			})
			id := a.newID()
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
			t.Run("ReadingUnpublishedEventsShouldReturnEventsWithExpectedData", func(t *testing.T) {
				err := a.eventStore.AddAggregate(context.Background(), testAggregateType, id, append(initialEvents, appendedEvents...))
				require.NoError(t, err)

				events, err := a.outbox.UnpublishedEvents(context.Background())
				_ = assert.NoError(t, err) &&
					assert.Len(t, events, 4) &&
					assert.Equal(t, expectedUnpublishedEvents, events)
			})
			t.Run("ReadingUnpublishedEventsAfterEverythingIsPublishedShouldReturnEmptyEvents", func(t *testing.T) {
				err := a.outbox.MarkEventsAsPublised(context.Background(), expectedUnpublishedEvents)
				_ = assert.NoError(t, err)

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
		assert.Equal(t, err, esrc.ErrOptimisticConcurrency)
}
