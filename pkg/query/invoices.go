package query

import "github.com/pperaltaisern/financing/pkg/financing"

type InvoiceQueries interface {
	ByID(financing.ID) (Invoice, error)
	All() ([]Invoice, error)
}

type Invoice struct {
	ID          financing.ID `gorm:"primaryKey"`
	IssuerID    financing.ID
	AskingPrice financing.Money
	WinningBid  *Bid `gorm:"embedded;embeddedPrefix:winning_bid_"`
	Status      InvoiceStatus
}

type Bid struct {
	InvestorID financing.ID `gorm:"primaryKey"`
	Amount     financing.Money
}

type InvoiceStatus string

const (
	InvoiceStatusAvailable InvoiceStatus = "Available"
	InvoiceStatusFinanced  InvoiceStatus = "Financed"
	InvoiceStatusApproved  InvoiceStatus = "Approved"
	InvoiceStatusReversed  InvoiceStatus = "Reversed"
)
