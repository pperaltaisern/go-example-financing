package financing

import "ledger/internal/es"

type Issuer struct {
	es.Aggregate
	id ID
}

func NewIssuer(id ID) *Issuer {
	inv := &Issuer{}
	inv.Aggregate = es.NewAggregate(inv.onEvent)

	e := NewIssuerCreatedEvent(id)
	inv.Raise(e)
	return inv
}

func (iss *Issuer) ID() ID {
	return iss.id
}

func NewIssuerFromEvents(events []es.Event) *Issuer {
	inv := &Issuer{}
	inv.Aggregate = es.NewAggregate(inv.onEvent)

	inv.Replay(events)

	return inv
}

func (iss *Issuer) onEvent(event es.Event) {
	switch e := event.(type) {
	case IssuerCreatedEvent:
		iss.id = e.IssuerID
	}
}
