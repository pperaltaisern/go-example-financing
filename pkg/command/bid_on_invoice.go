package command

import "ledger/pkg/financing"

type BidOnInvoice struct {
	InvoiceID  financing.ID
	InvestorID financing.ID
	BidAmount  financing.Money
}
