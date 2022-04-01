package esrc

import "context"

type MockEventStore struct {
	AddAggregateFn      func(context.Context, AggregateType, ID, []RawEvent) error
	AppendEventsFn      func(ctx context.Context, t AggregateType, id ID, fromVersion int, events []RawEvent) error
	AddSnapshotFn       func(context.Context, AggregateType, ID, RawSnapshot) error
	LatestSnapshotFn    func(context.Context, AggregateType, ID) (*RawSnapshot, error)
	EventsFn            func(context.Context, AggregateType, ID, int) ([]RawEvent, error)
	ContainsAggregateFn func(context.Context, AggregateType, ID) (bool, error)
}

var _ EventStore = (*MockEventStore)(nil)

func (m *MockEventStore) AddAggregate(ctx context.Context, t AggregateType, id ID, re []RawEvent) error {
	return m.AddAggregateFn(ctx, t, id, re)
}

func (m *MockEventStore) AppendEvents(ctx context.Context, t AggregateType, id ID, fromVersion int, events []RawEvent) error {
	return m.AppendEventsFn(ctx, t, id, fromVersion, events)
}

func (m *MockEventStore) AddSnapshot(ctx context.Context, t AggregateType, id ID, snapshot RawSnapshot) error {
	return m.AddSnapshotFn(ctx, t, id, snapshot)
}

func (m *MockEventStore) LatestSnapshot(ctx context.Context, t AggregateType, id ID) (*RawSnapshot, error) {
	return m.LatestSnapshotFn(ctx, t, id)
}

func (m *MockEventStore) Events(ctx context.Context, t AggregateType, id ID, fromEventNumber int) ([]RawEvent, error) {
	return m.EventsFn(ctx, t, id, fromEventNumber)
}
func (m *MockEventStore) ContainsAggregate(ctx context.Context, t AggregateType, id ID) (bool, error) {
	return m.ContainsAggregateFn(ctx, t, id)
}

type MockEvent struct {
	EventNameFn func() string
}

func (m *MockEvent) EventName() string {
	return m.EventNameFn()
}

type MockEventFactory struct {
	CreateEmptyEventFn func(name string) (Event, error)
}

func (m *MockEventFactory) CreateEmptyEvent(name string) (Event, error) {
	return m.CreateEmptyEventFn(name)
}

type MockEventMarshaler struct {
	MarshalEventFn   func(Event) ([]byte, error)
	UnmarshalEventFn func([]byte, Event) error
}

func (m *MockEventMarshaler) MarshalEvent(e Event) ([]byte, error) {
	return m.MarshalEventFn(e)
}

func (m *MockEventMarshaler) UnmarshalEvent(b []byte, e Event) error {
	return m.UnmarshalEventFn(b, e)
}
