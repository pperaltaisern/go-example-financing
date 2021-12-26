package command

import "ledger/pkg/financing"

type CreateInvestor struct {
	ID      financing.ID
	Balance financing.Money
}
