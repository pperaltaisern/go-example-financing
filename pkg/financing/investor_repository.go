package financing

import (
	"context"
	"fmt"

	"github.com/pperaltaisern/financing/internal/esrc"
)

type InvestorRepository interface {
	Update(context.Context, ID, UpdateInvestor) error
	Add(context.Context, *Investor) error
}

type UpdateInvestor func(inv *Investor) error

type investorRepository struct {
	r *esrc.Repository
}

func NewInvestorRepository(es esrc.EventStore) InvestorRepository {
	return investorRepository{
		r: esrc.NewRepository("investor", es, investorEventsFactory{}, esrc.JSONEventMarshaler{}),
	}
}

type investorEventsFactory struct{}

func (investorEventsFactory) CreateEmptyEvent(name string) (esrc.Event, error) {
	var e esrc.Event
	switch name {
	case "InvestorCreatedEvent":
		e = &InvestorCreatedEvent{}
	case "InvestorFundsAddedEvent":
		e = &InvestorFundsAddedEvent{}
	case "BidOnInvoicePlacedEvent":
		e = &BidOnInvoicePlacedEvent{}
	case "InvestorFundsReleasedEvent":
		e = &InvestorFundsReleasedEvent{}
	case "InvestorFundsCommittedEvent":
		e = &InvestorFundsCommittedEvent{}
	default:
		return nil, fmt.Errorf("unkown event name: %s", name)
	}
	return e, nil
}

func (r investorRepository) Update(ctx context.Context, id ID, update UpdateInvestor) error {
	inv, err := r.byID(ctx, id)
	if err != nil {
		return err
	}

	err = update(inv)
	if err != nil {
		return err
	}

	return r.r.Update(ctx, id, inv.aggregate.Version(), inv.aggregate.Events())
}

func (r investorRepository) byID(ctx context.Context, id ID) (*Investor, error) {
	events, err := r.r.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return newInvestorFromEvents(events), nil
}

func (r investorRepository) Add(ctx context.Context, inv *Investor) error {
	return r.r.Add(ctx, inv.id, inv.aggregate.Events())
}
