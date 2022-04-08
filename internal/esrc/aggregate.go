package esrc

type AggregateFactory interface {
	NewAggregateFromSnapshotAndEvents(RawSnapshot, []Event) (Aggregate, error)
	NewAggregateFromEvents([]Event) (Aggregate, error)
}

type Aggregate interface {
	InitialVersion() int
	Changes() []Event
	Snapshot() (RawSnapshot, error)
}

// EventRaiserAggregate is a helper struct that can be embbeded in real aggregates,
// handles the execution of raised events and keeps track of versioning
type EventRaiserAggregate struct {
	initialVersion int
	changes        []Event
	onEvent        func(Event)
}

func NewEventRaiserAggregate(onEvent func(Event)) EventRaiserAggregate {
	return EventRaiserAggregate{
		onEvent: onEvent,
	}
}

func NewEventRaiserAggregateFromEvents(initialVersion int, events []Event, onEvent func(Event)) EventRaiserAggregate {
	a := NewEventRaiserAggregate(onEvent)
	a.initialVersion = initialVersion
	a.replay(events)
	return a
}

func (a EventRaiserAggregate) Changes() []Event {
	return a.changes
}

func (a EventRaiserAggregate) InitialVersion() int {
	return a.initialVersion
}

func (a *EventRaiserAggregate) replay(events []Event) {
	a.initialVersion += len(events)
	for _, e := range events {
		a.onEvent(e)
	}
}

func (a *EventRaiserAggregate) Raise(e Event) {
	a.changes = append(a.changes, e)
	a.onEvent(e)
}
