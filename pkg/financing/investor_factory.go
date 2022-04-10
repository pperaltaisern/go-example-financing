package financing

import (
	"encoding/json"

	"github.com/pperaltaisern/financing/internal/esrc"
)

type investorFactory struct{}

var _ esrc.AggregateFactory[*Investor] = (*investorFactory)(nil)

func (investorFactory) NewAggregateFromSnapshotAndEvents(snapshot esrc.RawSnapshot, events []esrc.Event) (*Investor, error) {
	var invSnapshot investorSnapshot
	err := json.Unmarshal(snapshot.Data, &invSnapshot)
	if err != nil {
		return nil, err
	}

	inv := &Investor{
		id:       invSnapshot.ID,
		balance:  invSnapshot.Balance,
		reserved: invSnapshot.Reserved,
	}
	inv.EventRaiserAggregate = esrc.NewEventRaiserAggregateFromEvents(snapshot.Version, events, inv.onEvent)
	return inv, nil
}

func (investorFactory) NewAggregateFromEvents(events []esrc.Event) (*Investor, error) {
	inv := &Investor{}
	inv.EventRaiserAggregate = esrc.NewEventRaiserAggregateFromEvents(0, events, inv.onEvent)
	return inv, nil
}
