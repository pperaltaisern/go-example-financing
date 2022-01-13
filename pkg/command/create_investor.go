package command

import "github.com/pperaltaisern/financing/pkg/financing"

type CreateInvestor struct {
	ID      financing.ID
	Balance financing.Money
}
