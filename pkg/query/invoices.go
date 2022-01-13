package query

import "github.com/pperaltaisern/financing/pkg/financing"

type InvoiceQueries interface {
	ByID(financing.ID) Invoice
	All() []Invoice
}

type Invoice struct {
	ID          financing.ID
	IssuerID    financing.ID
	AskingPrice financing.Money
	WinningBid  Bid
	Status      string
}

type Bid struct {
	InvestorID financing.ID
	Amount     financing.Money
}
