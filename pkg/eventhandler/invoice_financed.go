package eventhandler

import (
	"context"

	"github.com/pperaltaisern/financing/pkg/financing"
)

type InvoiceFinancedHandler struct {
	investors financing.InvestorRepository
}

func NewInvoiceFinancedHandler(r financing.InvestorRepository) *InvoiceFinancedHandler {
	return &InvoiceFinancedHandler{
		investors: r,
	}
}
func (h *InvoiceFinancedHandler) Handle(ctx context.Context, event *financing.InvoiceFinancedEvent) error {
	unplacedBid := event.Bid.Amount - event.AskingPrice
	if unplacedBid <= 0 {
		return nil
	}

	return h.investors.Update(ctx, event.Bid.InvestorID, func(investor *financing.Investor) error {
		return investor.ReleaseFunds(unplacedBid)
	})
}
