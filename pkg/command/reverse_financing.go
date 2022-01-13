package command

import "github.com/pperaltaisern/financing/pkg/financing"

type ReverseFinancing struct {
	InvoiceID financing.ID
}

func NewReverseFinancing(invoiceID financing.ID) *ReverseFinancing {
	return &ReverseFinancing{
		InvoiceID: invoiceID,
	}
}
