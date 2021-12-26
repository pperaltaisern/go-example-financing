package financing

import "ledger/internal/esrc"

type Issuer struct {
	esrc.Aggregate
	id ID
}

func NewIssuer(id ID) *Issuer {
	inv := &Issuer{}
	inv.Aggregate = esrc.NewAggregate(inv.onEvent)

	e := NewIssuerCreatedEvent(id)
	inv.Raise(e)
	return inv
}

func (iss *Issuer) ID() ID {
	return iss.id
}

func NewIssuerFromEvents(events []esrc.Event) *Issuer {
	inv := &Issuer{}
	inv.Aggregate = esrc.NewAggregate(inv.onEvent)

	inv.Replay(events)

	return inv
}

func (iss *Issuer) onEvent(event esrc.Event) {
	switch e := event.(type) {
	case IssuerCreatedEvent:
		iss.id = e.IssuerID
	}
}
