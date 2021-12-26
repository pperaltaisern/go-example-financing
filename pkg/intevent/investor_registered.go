package intevent

import (
	"ledger/pkg/financing"
)

type InvestorRegistered struct {
	ID      financing.ID
	Name    string
	Balance financing.Money
}
