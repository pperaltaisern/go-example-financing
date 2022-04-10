package eventhandler

import (
	"context"

	"github.com/pperaltaisern/financing/pkg/financing"
)

type BidOnInvoicePlacedHandler struct {
	invoices financing.InvoiceRepository
}

func NewBidOnInvoicePlacedHandler(r financing.InvoiceRepository) *BidOnInvoicePlacedHandler {
	return &BidOnInvoicePlacedHandler{
		invoices: r,
	}
}

func (h *BidOnInvoicePlacedHandler) Handle(ctx context.Context, event *financing.BidOnInvoicePlacedEvent) error {
	return h.invoices.Update(ctx, event.InvoiceID, func(invoice *financing.Invoice) error {
		bid := financing.NewBid(event.InvestorID, event.BidAmount)
		invoice.ProcessBid(bid)
		return nil
	})
}
