package query

import "github.com/pperaltaisern/financing/pkg/financing"

type InvestorQueries interface {
	All() []Investor
}

type Investor struct {
	ID       financing.ID
	Balance  financing.Money
	Reserved financing.Money
}
