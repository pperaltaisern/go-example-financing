package command

import (
	"context"

	"github.com/pperaltaisern/financing/pkg/financing"
)

type ReverseFinancing struct {
	InvoiceID financing.ID
}

func NewReverseFinancing(invoiceID financing.ID) *ReverseFinancing {
	return &ReverseFinancing{
		InvoiceID: invoiceID,
	}
}

type ReverseFinancingHandler struct {
	invoices financing.InvoiceRepository
}

func NewReverseFinancingHandler(r financing.InvoiceRepository) *ReverseFinancingHandler {
	return &ReverseFinancingHandler{
		invoices: r,
	}
}

func (h *ReverseFinancingHandler) Handle(ctx context.Context, cmd *ReverseFinancing) error {
	return h.invoices.Update(ctx, cmd.InvoiceID, func(invoice *financing.Invoice) error {
		invoice.ReverseFinancing()
		return nil
	})
}
