package query

import "github.com/pperaltaisern/financing/pkg/financing"

type InvestorQueries interface {
	All() ([]Investor, error)
}

type Investor struct {
	ID        financing.ID
	Balance   financing.Money
	Reserved  financing.Money
	Committed financing.Money
}
