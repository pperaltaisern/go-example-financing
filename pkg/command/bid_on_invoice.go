package command

import "github.com/pperaltaisern/financing/pkg/financing"

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
