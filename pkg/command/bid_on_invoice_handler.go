package command

import (
	"context"
	"ledger/pkg/financing"

	"github.com/ThreeDotsLabs/watermill/components/cqrs"
)

type BidOnInvoiceHandler struct {
	investors financing.InvestorRepository
}

func NewBidOnInvoiceHandler(r financing.InvestorRepository) *BidOnInvoiceHandler {
	return &BidOnInvoiceHandler{
		investors: r,
	}
}

var _ cqrs.CommandHandler = (*BidOnInvoiceHandler)(nil)

func (h *BidOnInvoiceHandler) HandlerName() string {
	return "BidOnInvoiceHandler"
}

func (h *BidOnInvoiceHandler) NewCommand() interface{} {
	return &BidOnInvoice{}
}

func (h *BidOnInvoiceHandler) Handle(ctx context.Context, c interface{}) error {
	cmd := c.(*BidOnInvoice)

	return h.investors.Update(ctx, cmd.InvestorID, func(investor *financing.Investor) error {
		return investor.BidOnInvoice(cmd.InvoiceID, cmd.BidAmount)
	})
}
