package financing

import "github.com/pperaltaisern/financing/internal/esrc"

type Issuer struct {
	esrc.EventRaiserAggregate
	id ID
}

func NewIssuer(id ID) *Issuer {
	inv := &Issuer{}
	inv.EventRaiserAggregate = esrc.NewEventRaiserAggregate(inv.onEvent)

	e := NewIssuerCreatedEvent(id)
	inv.Raise(e)
	return inv
}

var _ esrc.Aggregate = (*Issuer)(nil)

func (iss *Issuer) ID() esrc.ID {
	return iss.id
}

func (iss *Issuer) onEvent(event esrc.Event) {
	switch e := event.(type) {
	case *IssuerCreatedEvent:
		iss.id = e.IssuerID
	}
}

func (iss *Issuer) Snapshot() ([]byte, error) {
	return nil, nil
}
