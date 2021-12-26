package eventhandler

import (
	"context"
	"ledger/pkg/financing"

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

	investor, err := h.investors.ByID(ctx, event.Bid.IvestorID)
	if err != nil {
		return err
	}

	investor.ReleaseFunds(event.Bid.Amount)

	h.investors.Update(ctx, investor)
	return nil
}
