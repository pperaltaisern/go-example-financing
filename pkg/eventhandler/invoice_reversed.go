package eventhandler

import (
	"context"

	"github.com/pperaltaisern/financing/pkg/financing"
)

type InvoiceReversedHandler struct {
	investors financing.InvestorRepository
}

func NewInvoiceReversedHandler(r financing.InvestorRepository) *InvoiceReversedHandler {
	return &InvoiceReversedHandler{
		investors: r,
	}
}

func (h *InvoiceReversedHandler) Handle(ctx context.Context, event *financing.InvoiceReversedEvent) error {
	return h.investors.Update(ctx, event.Bid.InvestorID, func(investor *financing.Investor) error {
		return investor.ReleaseFunds(event.SoldPrice)
	})
}
