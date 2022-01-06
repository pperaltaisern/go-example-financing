package financing

import (
	"context"
	"encoding/json"
	"fmt"
	"ledger/internal/esrc"
)

type InvestorRepository interface {
	ByID(context.Context, ID) (*Investor, error)
	Update(context.Context, *Investor) error
	Add(context.Context, *Investor) error
}

type investorRepository struct {
	es esrc.EventStore
}

func (i *investorRepository) byID(ctx context.Context, id ID) (*Investor, error) {
	rawEvents, err := i.es.Load(ctx, esrc.ID(id))
	if err != nil {
		return nil, err
	}

	events := make([]esrc.Event, len(rawEvents))
	for i, raw := range rawEvents {
		var e esrc.Event
		switch raw.Name {
		case "InvestorCreatedEvent":
			e = InvestorCreatedEvent{}
		case "InvestorFundsAddedEvent":
			e = InvestorFundsAddedEvent{}
		case "BidOnInvoicePlacedEvent":
			e = BidOnInvoicePlacedEvent{}
		case "InvestorFundsReleasedEvent":
			e = InvestorFundsReleasedEvent{}
		default:
			return nil, fmt.Errorf("unkown event name: %s", raw.Name)
		case "":
		}
		err = json.Unmarshal(raw.Data, e)
		if err != nil {
			return nil, err
		}
		events[i] = e
	}

	return newInvestorFromEvents(id, events), nil
}

func (i *investorRepository) Update(ctx context.Context, inv *Investor) error {
	rawEvents := make([]esrc.RawEvent, len(inv.aggregate.Events()))
	for i, e := range inv.aggregate.Events() {
		b, err := json.Marshal(e)
		if err != nil {
			return err
		}
		rawEvents[i] = esrc.RawEvent{Name: e.EventName(), Data: b}
	}

	return i.es.AppendEvents(ctx, inv.id, inv.aggregate.Version(), rawEvents)
}

func (i *investorRepository) Add(ctx context.Context, inv *Investor) error {
	rawEvents := make([]esrc.RawEvent, len(inv.aggregate.Events()))
	for i, e := range inv.aggregate.Events() {
		b, err := json.Marshal(e)
		if err != nil {
			return err
		}
		rawEvents[i] = esrc.RawEvent{Name: e.EventName(), Data: b}
	}

	const aggregateType esrc.AggregateType = "investor"
	return i.es.Create(ctx, aggregateType, inv.id, rawEvents)
}
