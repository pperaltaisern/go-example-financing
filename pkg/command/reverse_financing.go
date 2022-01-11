package command

import "ledger/pkg/financing"

type ReverseFinancing struct {
	InvoiceID financing.ID
}

func NewReverseFinancing(invoiceID financing.ID) *ReverseFinancing {
	return &ReverseFinancing{
		InvoiceID: invoiceID,
	}
}
