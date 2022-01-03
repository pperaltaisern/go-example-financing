package command

import "ledger/pkg/financing"

type SellInvoice struct {
	InvoiceID   financing.ID
	IssuerID    financing.ID
	AskingPrice financing.Money
}
