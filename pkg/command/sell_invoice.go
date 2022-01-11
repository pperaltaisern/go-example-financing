package command

import "ledger/pkg/financing"

type SellInvoice struct {
	InvoiceID   financing.ID
	IssuerID    financing.ID
	AskingPrice financing.Money
}

func NewSellInvoice(invoiceID, issuerID financing.ID, askingPrice financing.Money) *SellInvoice {
	return &SellInvoice{
		InvoiceID:   invoiceID,
		IssuerID:    issuerID,
		AskingPrice: askingPrice,
	}
}
