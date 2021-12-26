package es

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

func (a Aggregate) Events() []Event {
	return a.changes
}

func (a Aggregate) Version() int {
	return a.version
}

func (a *Aggregate) Raise(e Event) {
	a.changes = append(a.changes, e)
	a.onEvent(e)
}

func (a *Aggregate) Replay(events []Event) {
	a.version = len(events)
	for _, e := range events {
		a.onEvent(e)
	}
}
