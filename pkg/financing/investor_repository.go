package financing

import (
	"context"

	"github.com/pperaltaisern/esrc"
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
		r: esrc.NewRepository[*Investor](es, investorFactory{}, investorEventsFactory{}, opts...),
	}
}

func (r investorRepository) Update(ctx context.Context, id ID, update UpdateInvestor) error {
	return r.r.UpdateByID(ctx, id, update)
}

func (r investorRepository) Add(ctx context.Context, inv *Investor) error {
	return r.r.Add(ctx, inv)
}
