package esrc

import "context"

type MockEventStore struct {
	LoadFn         func(context.Context, AggregateType, ID) ([]RawEvent, error)
	ContainsFn     func(context.Context, AggregateType, ID) (bool, error)
	CreateFn       func(context.Context, AggregateType, ID, []RawEvent) error
	AppendEventsFn func(ctx context.Context, t AggregateType, id ID, fromVersion int, events []RawEvent) error
}

var _ EventStore = (*MockEventStore)(nil)

func (m *MockEventStore) Load(ctx context.Context, t AggregateType, id ID) ([]RawEvent, error) {
	return m.LoadFn(ctx, t, id)
}
func (m *MockEventStore) Contains(ctx context.Context, t AggregateType, id ID) (bool, error) {
	return m.ContainsFn(ctx, t, id)
}
func (m *MockEventStore) Create(ctx context.Context, t AggregateType, id ID, re []RawEvent) error {
	return m.CreateFn(ctx, t, id, re)
}
func (m *MockEventStore) AppendEvents(ctx context.Context, t AggregateType, id ID, fromVersion int, events []RawEvent) error {
	return m.AppendEventsFn(ctx, t, id, fromVersion, events)
}
