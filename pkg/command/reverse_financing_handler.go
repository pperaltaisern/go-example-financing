package command

import (
	"context"
	"ledger/pkg/financing"
)

type ReverseFinancingHandler struct {
	invoices financing.InvoiceRepository
}

func NewReverseFinancingHandler(r financing.InvoiceRepository) *ReverseFinancingHandler {
	return &ReverseFinancingHandler{
		invoices: r,
	}
}

func (h *ReverseFinancingHandler) HandlerName() string {
	return "ReverseFinancingHandler"
}

func (h *ReverseFinancingHandler) NewCommand() interface{} {
	return &ReverseFinancing{}
}

func (h *ReverseFinancingHandler) Handle(ctx context.Context, c interface{}) error {
	cmd := c.(*ReverseFinancing)

	return h.invoices.Update(ctx, cmd.InvoiceID, func(invoice *financing.Invoice) error {
		invoice.ReverseFinancing()
		return nil
	})
}
