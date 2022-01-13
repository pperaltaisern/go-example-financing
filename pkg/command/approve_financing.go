package command

import "github.com/pperaltaisern/financing/pkg/financing"

type ApproveFinancing struct {
	InvoiceID financing.ID
}

func NewApproveFinancing(invoiceID financing.ID) *ApproveFinancing {
	return &ApproveFinancing{
		InvoiceID: invoiceID,
	}
}
