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

	invoice, err := h.invoices.ByID(ctx, cmd.InvoiceID)
	if err != nil {
		return err
	}

	invoice.ReverseFinancing()

	err = h.invoices.Update(ctx, invoice)
	if err != nil {
		return err
	}

	return nil
}
