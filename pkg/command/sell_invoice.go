package command

import "ledger/pkg/financing"

type SellInvoice struct {
	IssuerID    financing.ID
	AskingPrice financing.Money
}
