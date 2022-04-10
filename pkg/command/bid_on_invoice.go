package command

import (
	"context"

	"github.com/pperaltaisern/financing/pkg/financing"
)

type BidOnInvoice struct {
	InvestorID financing.ID
	InvoiceID  financing.ID
	BidAmount  financing.Money
}

func NewBidOnInvoice(investorID, invoiceID financing.ID, bidAmount financing.Money) *BidOnInvoice {
	return &BidOnInvoice{
		InvestorID: investorID,
		InvoiceID:  invoiceID,
		BidAmount:  bidAmount,
	}
}

type BidOnInvoiceHandler struct {
	investors financing.InvestorRepository
}

func NewBidOnInvoiceHandler(r financing.InvestorRepository) *BidOnInvoiceHandler {
	return &BidOnInvoiceHandler{
		investors: r,
	}
}

func (h *BidOnInvoiceHandler) Handle(ctx context.Context, cmd *BidOnInvoice) error {
	return h.investors.Update(ctx, cmd.InvestorID, func(investor *financing.Investor) error {
		return investor.BidOnInvoice(cmd.InvoiceID, cmd.BidAmount)
	})
}
