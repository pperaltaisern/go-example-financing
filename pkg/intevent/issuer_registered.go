package intevent

import (
	"ledger/pkg/financing"
)

type IssuerRegistered struct {
	ID     financing.ID
	Name   string
	Amount financing.Money
}
