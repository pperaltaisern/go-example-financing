package esrc

// Aggregate is a helper base struct that can be embbeded in real aggregates
// that handles the execution of events raised and keeps trak of the versioning
type Aggregate struct {
	changes []Event
	version int

	onEvent func(Event)
}

func NewAggregate(onEvent func(Event)) Aggregate {
	return Aggregate{
		onEvent: onEvent,
	}
}

func NewAggregateFromEvents(events []Event, onEvent func(Event)) Aggregate {
	a := NewAggregate(onEvent)
	a.replay(events)
	return a
}

func (a *Aggregate) replay(events []Event) {
	a.version = len(events)
	for _, e := range events {
		a.onEvent(e)
	}
}

func (a *Aggregate) Raise(e Event) {
	a.changes = append(a.changes, e)
	a.onEvent(e)
}

func (a Aggregate) Events() []Event {
	return a.changes
}

func (a Aggregate) Version() int {
	return a.version
}
