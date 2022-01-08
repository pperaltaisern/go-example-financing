package esrc

import "context"

type MockEventStore struct {
	LoadFn         func(context.Context, ID) ([]RawEvent, error)
	ContainsFn     func(context.Context, ID) (bool, error)
	CreateFn       func(context.Context, AggregateType, ID, []RawEvent) error
	AppendEventsFn func(ctx context.Context, id ID, fromVersion int, events []RawEvent) error
}

var _ EventStore = (*MockEventStore)(nil)

func (m *MockEventStore) Load(ctx context.Context, id ID) ([]RawEvent, error) {
	return m.LoadFn(ctx, id)
}
func (m *MockEventStore) Contains(ctx context.Context, id ID) (bool, error) {
	return m.ContainsFn(ctx, id)
}
func (m *MockEventStore) Create(ctx context.Context, t AggregateType, id ID, re []RawEvent) error {
	return m.CreateFn(ctx, t, id, re)
}
func (m *MockEventStore) AppendEvents(ctx context.Context, id ID, fromVersion int, events []RawEvent) error {
	return m.AppendEventsFn(ctx, id, fromVersion, events)
}
