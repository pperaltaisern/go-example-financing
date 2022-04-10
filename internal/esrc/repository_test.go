package esrc

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRepository(t *testing.T) {
	aggregateType := AggregateType("MockAggregate")
	ctx := context.Background()
	eventStoreErr := ErrAggregateNotFound
	foundEvent := &testEvent{}
	rawFoundEvent := RawEvent{Name: foundEvent.EventName(), Data: []byte("{}")}

	t.Run("FindByIdShould", func(t *testing.T) {
		dontFindSnapshot := func(_ context.Context, _ AggregateType, _ ID) (*RawSnapshot, error) {
			return nil, nil
		}
		findAnEvent := func(_ context.Context, _ AggregateType, _ ID, _ int) ([]RawEvent, error) {
			return []RawEvent{rawFoundEvent}, nil
		}

		t.Run("ScaleTheErrorsOfLatestSnapshot", func(t *testing.T) {
			eventStore := &MockEventStore{
				LatestSnapshotFn: func(_ context.Context, _ AggregateType, _ ID) (*RawSnapshot, error) {
					return nil, eventStoreErr
				},
			}

			repository := NewRepository[*MockAggregate](eventStore, nil, nil)
			aggregate, err := repository.FindByID(ctx, eventStore)
			assert.Nil(t, aggregate)
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

			repository := NewRepository[*MockAggregate](eventStore, nil, nil)
			aggregate, err := repository.FindByID(ctx, eventStore)
			assert.Nil(t, aggregate)
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

			repository := NewRepository[*MockAggregate](eventStore, nil, nil)
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

			repository := NewRepository[*MockAggregate](eventStore, nil, nil)
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
					assert.Equal(t, rawFoundEvent.Name, name)
					return nil, factoryErr
				},
			}

			repository := NewRepository[*MockAggregate](eventStore, nil, eventFactory)
			aggregate, err := repository.FindByID(ctx, nil)
			assert.Nil(t, aggregate)
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
					assert.Equal(t, string(rawFoundEvent.Data), string(b))
					return marshalerErr
				},
			}

			repository := NewRepository(eventStore, nil, eventFactory, RepositoryWithEventsMarshaler[*MockAggregate](eventMarshaler))
			aggregate, err := repository.FindByID(ctx, eventStore)
			assert.Nil(t, aggregate)
			assert.Equal(t, marshalerErr, err)
		})

		t.Run("CreateAggregateWithFoundSnapshotAndEvents", func(t *testing.T) {
			foundSnapshot := RawSnapshot{Version: 5}

			expectedAggregate := &MockAggregate{IDFn: func() ID { return "id" }}
			aggregateFactory := &MockAggregateFactory[*MockAggregate]{
				NewAggregateFromSnapshotAndEventsFn: func(rs RawSnapshot, e []Event) (*MockAggregate, error) {
					require.Equal(t, foundSnapshot, rs)
					require.Equal(t, []Event{foundEvent}, e)
					return expectedAggregate, nil
				},
			}

			eventStore := &MockEventStore{
				LatestSnapshotFn: func(_ context.Context, _ AggregateType, _ ID) (*RawSnapshot, error) {
					return &foundSnapshot, nil
				},
				EventsFn: findAnEvent,
			}
			eventFactory := &MockEventFactory{
				CreateEmptyEventFn: func(name string) (Event, error) {
					return foundEvent, nil
				},
			}

			repository := NewRepository[*MockAggregate](eventStore, aggregateFactory, eventFactory)
			aggregate, err := repository.FindByID(ctx, eventStore)
			require.NoError(t, err)
			assert.Equal(t, expectedAggregate, aggregate)
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
			repository := NewRepository[*MockAggregate](eventStore, nil, nil)

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
			repository := NewRepository[*MockAggregate](eventStore, nil, nil)

			_, err := repository.Contains(ctx, "")
			require.Error(t, err)
		})
	})

	t.Run("UpdateShould", func(t *testing.T) {
		id := 10
		fromVersion := 9
		rawEvents := []RawEvent{rawFoundEvent}
		events := []Event{foundEvent}
		succesfulUpdate := func(_ context.Context, at AggregateType, thisID ID, thisFromVersion int, thisEvents []RawEvent) error {
			return nil
		}

		t.Run("CallAppendEventsAndScaleErrors", func(t *testing.T) {
			aggregate := &MockAggregate{
				IDFn:             func() ID { return id },
				InitialVersionFn: func() int { return fromVersion },
				ChangesFn:        func() []Event { return events },
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
			repository := NewRepository[*MockAggregate](eventStore, nil, nil)

			err := repository.Update(ctx, aggregate)
			require.Equal(t, eventStoreErr, err)
		})

		t.Run("NotCreateSnapshotIfEventLimitIsNotReached", func(t *testing.T) {
			aggregate := &MockAggregate{
				IDFn:             func() ID { return id },
				InitialVersionFn: func() int { return fromVersion },
				ChangesFn:        func() []Event { return events },
				SnapshotFn:       func() ([]byte, error) { panic("this shouldn't be called") },
			}

			eventStore := &MockEventStore{
				AppendEventsFn: succesfulUpdate,
			}
			repository := NewRepository(
				eventStore,
				nil,
				nil,
				RepositoryWithEventsPerSnapshot[*MockAggregate](200))

			err := repository.Update(ctx, aggregate)
			require.NoError(t, err)
		})

		t.Run("CreateSnapshotIfEventLimitIsReached", func(t *testing.T) {
			snapshotData := []byte("test snapshot")
			aggregate := &MockAggregate{
				IDFn:             func() ID { return id },
				InitialVersionFn: func() int { return fromVersion },
				ChangesFn:        func() []Event { return events },
				SnapshotFn:       func() ([]byte, error) { return snapshotData, nil },
			}

			addSnapshotHasBeenCalled := false
			eventStore := &MockEventStore{
				AppendEventsFn: succesfulUpdate,
				AddSnapshotFn: func(ctx context.Context, at AggregateType, thisID ID, rs RawSnapshot) error {
					require.Equal(t, aggregateType, at)
					require.Equal(t, id, thisID)
					require.Equal(t, RawSnapshot{Version: 10, Data: snapshotData}, rs)
					addSnapshotHasBeenCalled = true
					return nil
				},
			}
			repository := NewRepository(
				eventStore,
				nil,
				nil,
				RepositoryWithEventsPerSnapshot[*MockAggregate](10))

			err := repository.Update(ctx, aggregate)
			require.NoError(t, err)
			// adding snapshots is done in background
			time.Sleep(100 * time.Millisecond)
			require.True(t, addSnapshotHasBeenCalled)
		})
	})

	t.Run("AddShould", func(t *testing.T) {
		id := 10
		fromVersion := 9
		rawEvents := []RawEvent{rawFoundEvent}
		events := []Event{foundEvent}
		addAggregateSuccessfuly := func(_ context.Context, _ AggregateType, _ ID, _ []RawEvent) error {
			return nil
		}

		t.Run("CallAddAggregateAndScaleErrors", func(t *testing.T) {
			aggregate := &MockAggregate{
				IDFn:      func() ID { return id },
				ChangesFn: func() []Event { return events },
			}
			eventStore := &MockEventStore{
				AddAggregateFn: func(_ context.Context, at AggregateType, thisID ID, thisEvents []RawEvent) error {
					assert.Equal(t, aggregateType, at)
					assert.Equal(t, id, thisID)
					assert.Equal(t, rawEvents, thisEvents)
					return eventStoreErr
				},
			}
			repository := NewRepository[*MockAggregate](eventStore, nil, nil)

			err := repository.Add(ctx, aggregate)
			require.Equal(t, eventStoreErr, err)
		})

		t.Run("NotCreateSnapshotIfEventLimitIsNotReached", func(t *testing.T) {
			aggregate := &MockAggregate{
				IDFn:             func() ID { return id },
				InitialVersionFn: func() int { return fromVersion },
				ChangesFn:        func() []Event { return events },
				SnapshotFn:       func() ([]byte, error) { panic("this shouldn't be called") },
			}

			eventStore := &MockEventStore{
				AddAggregateFn: addAggregateSuccessfuly,
			}
			repository := NewRepository(
				eventStore,
				nil,
				nil,
				RepositoryWithEventsPerSnapshot[*MockAggregate](200))

			err := repository.Add(ctx, aggregate)
			require.NoError(t, err)
		})

		t.Run("CreateSnapshotIfEventLimitIsReached", func(t *testing.T) {
			snapshotData := []byte("test snapshot")
			aggregate := &MockAggregate{
				IDFn:             func() ID { return id },
				InitialVersionFn: func() int { return fromVersion },
				ChangesFn:        func() []Event { return events },
				SnapshotFn:       func() ([]byte, error) { return snapshotData, nil },
			}

			addSnapshotHasBeenCalled := false
			eventStore := &MockEventStore{
				AddAggregateFn: addAggregateSuccessfuly,
				AddSnapshotFn: func(ctx context.Context, at AggregateType, thisID ID, rs RawSnapshot) error {
					require.Equal(t, aggregateType, at)
					require.Equal(t, id, thisID)
					require.Equal(t, RawSnapshot{Version: 10, Data: snapshotData}, rs)
					addSnapshotHasBeenCalled = true
					return nil
				},
			}
			repository := NewRepository(
				eventStore,
				nil,
				nil,
				RepositoryWithEventsPerSnapshot[*MockAggregate](10))

			err := repository.Add(ctx, aggregate)
			require.NoError(t, err)
			// adding snapshots is done in background
			time.Sleep(100 * time.Millisecond)
			require.True(t, addSnapshotHasBeenCalled)
		})
	})
}

func TestRepositoryShouldDoSnapshotWhenEventLimitSetAndReached(t *testing.T) {
	testCases := map[string]struct {
		EventsPerSnapshot int
		InitialVersion    int
		NumOfChanges      int
		ShouldDoSnapshot  bool
	}{
		"EventLimitNotSet": {
			EventsPerSnapshot: 0,
			InitialVersion:    9,
			NumOfChanges:      10,
			ShouldDoSnapshot:  false,
		},
		"LimitNotSurpassed": {
			EventsPerSnapshot: 20,
			InitialVersion:    9,
			NumOfChanges:      10,
			ShouldDoSnapshot:  false,
		},
		"LimitReached": {
			EventsPerSnapshot: 20,
			InitialVersion:    10,
			NumOfChanges:      10,
			ShouldDoSnapshot:  true,
		},
		"LimitSurpassed": {
			EventsPerSnapshot: 20,
			InitialVersion:    10,
			NumOfChanges:      11,
			ShouldDoSnapshot:  true,
		},
		"LimitSurpassedByTwiceTheLimit": {
			EventsPerSnapshot: 20,
			InitialVersion:    19,
			NumOfChanges:      22,
			ShouldDoSnapshot:  true,
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			r := NewRepository(nil, nil, nil, RepositoryWithEventsPerSnapshot[*MockAggregate](testCase.EventsPerSnapshot))

			agg := &MockAggregate{
				InitialVersionFn: func() int { return testCase.InitialVersion },
				ChangesFn:        func() []Event { return make([]Event, testCase.NumOfChanges) },
			}
			assert.Equal(t, testCase.ShouldDoSnapshot, r.shouldDoSnapshot(agg))
		})
	}
}
