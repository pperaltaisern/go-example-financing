package esrc

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRepository(t *testing.T) {
	aggregateType := AggregateType("testType")
	ctx := context.Background()
	eventStoreErr := ErrAggregateNotFound

	t.Run("FindByIdShould", func(t *testing.T) {
		foundRawEvent := RawEvent{Name: "found event", Data: []byte("found event data")}
		dontFindSnapshot := func(_ context.Context, _ AggregateType, _ ID) (*RawSnapshot, error) {
			return nil, nil
		}
		findAnEvent := func(_ context.Context, _ AggregateType, _ ID, _ int) ([]RawEvent, error) {
			return []RawEvent{foundRawEvent}, nil
		}

		t.Run("ScaleTheErrorsOfLatestSnapshot", func(t *testing.T) {
			eventStore := &MockEventStore{
				LatestSnapshotFn: func(_ context.Context, _ AggregateType, _ ID) (*RawSnapshot, error) {
					return nil, eventStoreErr
				},
			}

			repository := NewRepository(aggregateType, eventStore, nil, nil)
			snapshot, events, err := repository.FindByID(ctx, eventStore)
			assert.Nil(t, snapshot)
			assert.Empty(t, events)
			assert.Equal(t, eventStoreErr, err)
		})

		t.Run("ScaleErrorsOnGetEvent", func(t *testing.T) {
			eventStore := &MockEventStore{
				LatestSnapshotFn: func(_ context.Context, _ AggregateType, _ ID) (*RawSnapshot, error) {
					return &RawSnapshot{}, nil
				},
				EventsFn: func(_ context.Context, _ AggregateType, _ ID, _ int) ([]RawEvent, error) {
					return nil, eventStoreErr
				},
			}

			repository := NewRepository(aggregateType, eventStore, nil, nil)
			snapshot, events, err := repository.FindByID(ctx, eventStore)
			assert.Nil(t, snapshot)
			assert.Empty(t, events)
			assert.Equal(t, eventStoreErr, err)
		})

		t.Run("RequestEventFromVersion0WhenNoSnapshotIsFound", func(t *testing.T) {
			capturedFromVersion := -1
			eventStore := &MockEventStore{
				LatestSnapshotFn: dontFindSnapshot,
				EventsFn: func(_ context.Context, _ AggregateType, _ ID, fromVersion int) ([]RawEvent, error) {
					capturedFromVersion = fromVersion
					return nil, eventStoreErr
				},
			}

			repository := NewRepository(aggregateType, eventStore, nil, nil)
			repository.FindByID(ctx, eventStore)
			assert.Equal(t, 0, capturedFromVersion)
		})
		t.Run("RequestEventFromSnapshotVersionWhenSnapshotIsFound", func(t *testing.T) {
			snapshotVersion := 10
			var capturedFromVersion int
			eventStore := &MockEventStore{
				LatestSnapshotFn: func(_ context.Context, _ AggregateType, _ ID) (*RawSnapshot, error) {
					return &RawSnapshot{Version: snapshotVersion}, nil
				},
				EventsFn: func(_ context.Context, _ AggregateType, _ ID, fromVersion int) ([]RawEvent, error) {
					capturedFromVersion = fromVersion
					return nil, eventStoreErr
				},
			}

			repository := NewRepository(aggregateType, eventStore, nil, nil)
			repository.FindByID(ctx, eventStore)
			assert.Equal(t, snapshotVersion+1, capturedFromVersion)
		})

		t.Run("ScaleErrorsOnEventFactory", func(t *testing.T) {
			factoryErr := errors.New("factory err")
			eventStore := &MockEventStore{
				LatestSnapshotFn: dontFindSnapshot,
				EventsFn:         findAnEvent,
			}
			eventFactory := &MockEventFactory{
				CreateEmptyEventFn: func(name string) (Event, error) {
					assert.Equal(t, foundRawEvent.Name, name)
					return nil, factoryErr
				},
			}

			repository := NewRepository(aggregateType, eventStore, eventFactory, nil)
			snapshot, events, err := repository.FindByID(ctx, eventStore)
			assert.Nil(t, snapshot)
			assert.Empty(t, events)
			assert.Equal(t, factoryErr, err)
		})

		t.Run("ScaleErrorsOnEventMarshaler", func(t *testing.T) {
			marshalerErr := errors.New("marshaler err")
			emptyEvent := &MockEvent{EventNameFn: func() string { return "empty event" }}

			eventStore := &MockEventStore{
				LatestSnapshotFn: dontFindSnapshot,
				EventsFn:         findAnEvent,
			}
			eventFactory := &MockEventFactory{
				CreateEmptyEventFn: func(name string) (Event, error) {
					return emptyEvent, nil
				},
			}
			eventMarshaler := &MockEventMarshaler{
				UnmarshalEventFn: func(b []byte, e Event) error {
					assert.Equal(t, string(foundRawEvent.Data), string(b))
					return marshalerErr
				},
			}

			repository := NewRepository(aggregateType, eventStore, eventFactory, eventMarshaler)
			snapshot, events, err := repository.FindByID(ctx, eventStore)
			assert.Nil(t, snapshot)
			assert.Empty(t, events)
			assert.Equal(t, marshalerErr, err)
		})

		t.Run("ReturnFoundSnapshotAndEvents", func(t *testing.T) {
			foundSnapshot := &RawSnapshot{Version: 5}
			foundEvent := &MockEvent{EventNameFn: func() string { return "empty event" }}

			eventStore := &MockEventStore{
				LatestSnapshotFn: func(_ context.Context, _ AggregateType, _ ID) (*RawSnapshot, error) {
					return foundSnapshot, nil
				},
				EventsFn: findAnEvent,
			}
			eventFactory := &MockEventFactory{
				CreateEmptyEventFn: func(name string) (Event, error) {
					return foundEvent, nil
				},
			}
			eventMarshaler := &MockEventMarshaler{
				UnmarshalEventFn: func(b []byte, e Event) error {
					return nil
				},
			}

			repository := NewRepository(aggregateType, eventStore, eventFactory, eventMarshaler)
			snapshot, events, err := repository.FindByID(ctx, eventStore)
			require.NoError(t, err)
			assert.Equal(t, foundSnapshot, snapshot)
			assert.Equal(t, []Event{foundEvent}, events)
		})
	})

	t.Run("ContainsShould", func(t *testing.T) {
		t.Run("CallEventStoresContainsAggregate", func(t *testing.T) {
			id := 10
			eventStore := &MockEventStore{
				ContainsAggregateFn: func(ctx context.Context, at AggregateType, thisID ID) (bool, error) {
					assert.Equal(t, aggregateType, at)
					assert.Equal(t, id, thisID)
					return true, nil
				},
			}
			repository := NewRepository(aggregateType, eventStore, nil, nil)

			found, err := repository.Contains(ctx, id)
			require.True(t, found)
			require.NoError(t, err)
		})
		t.Run("ScaleErrorsOnEventStoresContainsAggregate", func(t *testing.T) {
			eventStore := &MockEventStore{
				ContainsAggregateFn: func(ctx context.Context, at AggregateType, thisID ID) (bool, error) {
					return false, eventStoreErr
				},
			}
			repository := NewRepository(aggregateType, eventStore, nil, nil)

			_, err := repository.Contains(ctx, "")
			require.Error(t, err)
		})
	})

	t.Run("UpdateShould", func(t *testing.T) {
		rawEvents := []RawEvent{{Name: "event"}}
		events := []Event{&MockEvent{EventNameFn: func() string { return "event" }}}
		t.Run("CallAppendEventsAndScaleErrors", func(t *testing.T) {
			id := 10
			fromVersion := 9
			eventMarshaler := &MockEventMarshaler{
				MarshalEventFn: func(e Event) ([]byte, error) {
					return nil, nil
				},
			}
			eventStore := &MockEventStore{
				AppendEventsFn: func(_ context.Context, at AggregateType, thisID ID, thisFromVersion int, thisEvents []RawEvent) error {
					assert.Equal(t, aggregateType, at)
					assert.Equal(t, id, thisID)
					assert.Equal(t, fromVersion, thisFromVersion)
					assert.Equal(t, rawEvents, thisEvents)
					return eventStoreErr
				},
			}
			repository := NewRepository(aggregateType, eventStore, nil, eventMarshaler)

			err := repository.Update(ctx, id, fromVersion, events)
			require.Equal(t, eventStoreErr, err)
		})
	})

	t.Run("AddShould", func(t *testing.T) {
		rawEvents := []RawEvent{{Name: "event"}}
		events := []Event{&MockEvent{EventNameFn: func() string { return "event" }}}
		t.Run("CallAddAggregateAndScaleErrors", func(t *testing.T) {
			id := 10
			eventMarshaler := &MockEventMarshaler{
				MarshalEventFn: func(e Event) ([]byte, error) {
					return nil, nil
				},
			}
			eventStore := &MockEventStore{
				AddAggregateFn: func(_ context.Context, at AggregateType, thisID ID, thisEvents []RawEvent) error {
					assert.Equal(t, aggregateType, at)
					assert.Equal(t, id, thisID)
					assert.Equal(t, rawEvents, thisEvents)
					return eventStoreErr
				},
			}
			repository := NewRepository(aggregateType, eventStore, nil, eventMarshaler)

			err := repository.Add(ctx, id, events)
			require.Equal(t, eventStoreErr, err)
		})
	})
}
