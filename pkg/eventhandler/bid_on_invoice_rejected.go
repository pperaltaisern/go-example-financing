package eventhandler

import (
	"context"

	"github.com/pperaltaisern/financing/pkg/financing"
)

type BidOnInvoiceRejectedHandler struct {
	investors financing.InvestorRepository
}

func NewBidOnInvoiceRejectedHandler(r financing.InvestorRepository) *BidOnInvoiceRejectedHandler {
	return &BidOnInvoiceRejectedHandler{
		investors: r,
	}
}

func (h *BidOnInvoiceRejectedHandler) Handle(ctx context.Context, event *financing.BidOnInvoiceRejectedEvent) error {
	return h.investors.Update(ctx, event.Bid.InvestorID, func(investor *financing.Investor) error {
		return investor.ReleaseFunds(event.Bid.Amount)
	})
}
