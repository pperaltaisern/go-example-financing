package intevent

import (
	"github.com/pperaltaisern/financing/pkg/financing"
)

type InvestorRegistered struct {
	ID      financing.ID
	Name    string
	Balance financing.Money
}
