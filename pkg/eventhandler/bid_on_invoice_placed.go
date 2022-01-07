package eventhandler

import (
	"context"
	"ledger/pkg/financing"

	"github.com/ThreeDotsLabs/watermill/components/cqrs"
)

type BidOnInvoicePlacedHandler struct {
	invoices financing.InvoiceRepository
}

func NewBidOnInvoicePlacedHandler(r financing.InvoiceRepository) *BidOnInvoicePlacedHandler {
	return &BidOnInvoicePlacedHandler{
		invoices: r,
	}
}

var _ cqrs.EventHandler = (*BidOnInvoicePlacedHandler)(nil)

func (h *BidOnInvoicePlacedHandler) HandlerName() string {
	return "BidOnInvoicePlacedHandler"
}

func (h *BidOnInvoicePlacedHandler) NewEvent() interface{} {
	return &financing.BidOnInvoicePlacedEvent{}
}

func (h *BidOnInvoicePlacedHandler) Handle(ctx context.Context, e interface{}) error {
	event := e.(*financing.BidOnInvoicePlacedEvent)

	return h.invoices.Update(ctx, event.InvoiceID, func(invoice *financing.Invoice) error {
		invoice.ProcessBid(event.Bid)
		return nil
	})
}
