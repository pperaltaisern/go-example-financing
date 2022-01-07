package command

import (
	"context"
	"ledger/pkg/financing"

	"github.com/ThreeDotsLabs/watermill/components/cqrs"
)

type ApproveFinancingHandler struct {
	invoices financing.InvoiceRepository
}

func NewApproveFinancingHandler(r financing.InvoiceRepository) *ApproveFinancingHandler {
	return &ApproveFinancingHandler{
		invoices: r,
	}
}

var _ cqrs.CommandHandler = (*ApproveFinancingHandler)(nil)

func (h *ApproveFinancingHandler) HandlerName() string {
	return "ApproveFinancingHandler"
}

func (h *ApproveFinancingHandler) NewCommand() interface{} {
	return &ApproveFinancing{}
}

func (h *ApproveFinancingHandler) Handle(ctx context.Context, c interface{}) error {
	cmd := c.(*ApproveFinancing)

	return h.invoices.Update(ctx, cmd.InvoiceID, func(invoice *financing.Invoice) error {
		invoice.ApproveFinancing()
		return nil
	})
}
