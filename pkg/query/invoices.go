package query

import "github.com/pperaltaisern/financing/pkg/financing"

type InvoiceQueries interface {
	ByID(financing.ID) (Invoice, error)
	All() ([]Invoice, error)
}

type Invoice struct {
	ID          financing.ID    `gorm:"type:uuid;primary_key;"`
	IssuerID    financing.ID    `gorm:"type:uuid"`
	AskingPrice financing.Money `gorm:"type:float;"`
	WinningBid  *Bid            `gorm:"embedded;embeddedPrefix:winning_bid_"`
	Status      InvoiceStatus
}

type Bid struct {
	InvestorID financing.ID    `gorm:"type:uuid"`
	Amount     financing.Money `gorm:"type:float;"`
}

type InvoiceStatus string

const (
	InvoiceStatusAvailable InvoiceStatus = "Available"
	InvoiceStatusFinanced  InvoiceStatus = "Financed"
	InvoiceStatusApproved  InvoiceStatus = "Approved"
	InvoiceStatusReversed  InvoiceStatus = "Reversed"
)
