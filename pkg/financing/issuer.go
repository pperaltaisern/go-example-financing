package financing

import "github.com/pperaltaisern/financing/internal/esrc"

type Issuer struct {
	aggregate esrc.Aggregate
	id        ID
}

func NewIssuer(id ID) *Issuer {
	inv := &Issuer{}
	inv.aggregate = esrc.NewAggregate(inv.onEvent)

	e := NewIssuerCreatedEvent(id)
	inv.aggregate.Raise(e)
	return inv
}

func newIssuerFromEvents(events []esrc.Event) *Issuer {
	inv := &Issuer{}
	inv.aggregate = esrc.NewAggregateFromEvents(events, inv.onEvent)
	return inv
}

func (iss *Issuer) onEvent(event esrc.Event) {
	switch e := event.(type) {
	case *IssuerCreatedEvent:
		iss.id = e.IssuerID
	}
}
