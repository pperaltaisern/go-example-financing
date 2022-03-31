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
	WinningBid  *Bid            `gorm:"ForeignKey:InvoiceID"`
	Status      InvoiceStatus
}

type Bid struct {
	InvoiceID  financing.ID    `gorm:"type:uuid,primaryKey"`
	InvestorID financing.ID    `gorm:"type:uuid,primaryKey"`
	Amount     financing.Money `gorm:"type:float;"`
}

type InvoiceStatus string

const (
	InvoiceStatusAvailable InvoiceStatus = "Available"
	InvoiceStatusFinanced  InvoiceStatus = "Financed"
	InvoiceStatusApproved  InvoiceStatus = "Approved"
	InvoiceStatusReversed  InvoiceStatus = "Reversed"
)
