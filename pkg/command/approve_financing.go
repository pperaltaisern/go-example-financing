package command

import (
	"context"

	"github.com/pperaltaisern/financing/pkg/financing"
)

type ApproveFinancing struct {
	InvoiceID financing.ID
}

func NewApproveFinancing(invoiceID financing.ID) *ApproveFinancing {
	return &ApproveFinancing{
		InvoiceID: invoiceID,
	}
}

type ApproveFinancingHandler struct {
	invoices financing.InvoiceRepository
}

func NewApproveFinancingHandler(r financing.InvoiceRepository) *ApproveFinancingHandler {
	return &ApproveFinancingHandler{
		invoices: r,
	}
}
func (h *ApproveFinancingHandler) Handle(ctx context.Context, cmd *ApproveFinancing) error {
	return h.invoices.Update(ctx, cmd.InvoiceID, func(invoice *financing.Invoice) error {
		invoice.ApproveFinancing()
		return nil
	})
}
