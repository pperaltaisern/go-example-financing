package financing

import "context"

type InvestorRepository interface {
	ByID(context.Context, ID) (*Investor, error)
	Update(context.Context, *Investor) error
	Add(context.Context, *Investor) error
}
