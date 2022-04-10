package eventhandler

import (
	"context"

	"github.com/pperaltaisern/financing/pkg/financing"
)

type InvoiceApprovedHandler struct {
	investors financing.InvestorRepository
}

func NewInvoiceApprovedHandler(r financing.InvestorRepository) *InvoiceApprovedHandler {
	return &InvoiceApprovedHandler{
		investors: r,
	}
}

func (h *InvoiceApprovedHandler) Handle(ctx context.Context, event *financing.InvoiceApprovedEvent) error {
	return h.investors.Update(ctx, event.Bid.InvestorID, func(investor *financing.Investor) error {
		return investor.CommitFunds(event.SoldPrice)
	})
}
