package esrc

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
