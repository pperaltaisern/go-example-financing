package financing

import (
	"context"

	"github.com/pperaltaisern/financing/internal/esrc"
)

type InvestorRepository interface {
	Update(context.Context, ID, UpdateInvestor) error
	Add(context.Context, *Investor) error
}

type UpdateInvestor func(inv *Investor) error

type investorRepository struct {
	r *esrc.Repository[*Investor]
}

func NewInvestorRepository(es esrc.EventStore, opts ...esrc.RepositoryOption[*Investor]) InvestorRepository {
	return investorRepository{
		r: esrc.NewRepository[*Investor](
			es,
			investorFactory{},
			investorEventsFactory{},
			opts...),
	}
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

	return r.r.Update(ctx, inv)
}

func (r investorRepository) byID(ctx context.Context, id ID) (*Investor, error) {
	return r.r.FindByID(ctx, id)
}

func (r investorRepository) Add(ctx context.Context, inv *Investor) error {
	return r.r.Add(ctx, inv)
}
