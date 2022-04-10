package financing

import (
	"encoding/json"

	"github.com/pperaltaisern/financing/internal/esrc"
)

type invoiceFactory struct{}

var _ esrc.AggregateFactory[*Invoice] = (*invoiceFactory)(nil)

func (invoiceFactory) NewAggregateFromSnapshotAndEvents(snapshot esrc.RawSnapshot, events []esrc.Event) (*Invoice, error) {
	var invSnapshot invoiceSnapshot
	err := json.Unmarshal(snapshot.Data, &invSnapshot)
	if err != nil {
		return nil, err
	}

	inv := &Invoice{
		id:          invSnapshot.ID,
		issuerID:    invSnapshot.IssuerID,
		askingPrice: invSnapshot.AskingPrice,
		status:      invSnapshot.Status,
		winningBid:  invSnapshot.WinningBid,
	}
	inv.EventRaiserAggregate = esrc.NewEventRaiserAggregateFromEvents(snapshot.Version, events, inv.onEvent)
	return inv, nil
}

func (invoiceFactory) NewAggregateFromEvents(events []esrc.Event) (*Invoice, error) {
	inv := &Invoice{}
	inv.EventRaiserAggregate = esrc.NewEventRaiserAggregateFromEvents(0, events, inv.onEvent)
	return inv, nil
}
