package eventhandler

import (
	"context"

	"github.com/pperaltaisern/financing/pkg/financing"

	"github.com/ThreeDotsLabs/watermill/components/cqrs"
)

type BidOnInvoiceRejectedHandler struct {
	investors financing.InvestorRepository
}

func NewBidOnInvoiceRejectedHandler(r financing.InvestorRepository) *BidOnInvoiceRejectedHandler {
	return &BidOnInvoiceRejectedHandler{
		investors: r,
	}
}

var _ cqrs.EventHandler = (*BidOnInvoiceRejectedHandler)(nil)

func (h *BidOnInvoiceRejectedHandler) HandlerName() string {
	return "BidOnInvoiceRejectedHandler"
}

func (h *BidOnInvoiceRejectedHandler) NewEvent() interface{} {
	return &financing.BidOnInvoicePlacedEvent{}
}

func (h *BidOnInvoiceRejectedHandler) Handle(ctx context.Context, e interface{}) error {
	event := e.(*financing.BidOnInvoicePlacedEvent)

	return h.investors.Update(ctx, event.InvestorID, func(investor *financing.Investor) error {
		return investor.ReleaseFunds(event.BidAmount)
	})
}
