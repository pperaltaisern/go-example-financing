package command

import "ledger/pkg/financing"

type ApproveFinancing struct {
	InvoiceID financing.ID
}

func NewApproveFinancing(invoiceID financing.ID) *ApproveFinancing {
	return &ApproveFinancing{
		InvoiceID: invoiceID,
	}
}
