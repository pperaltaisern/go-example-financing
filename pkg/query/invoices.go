package query

import (
	"github.com/pperaltaisern/financing/pkg/financing"
	"gorm.io/gorm"
)

type InvoiceQueries interface {
	ByID(financing.ID) (Invoice, error)
	All() ([]Invoice, error)
}

type Invoice struct {
	ID          financing.ID    `gorm:"type:uuid;primary_key;"`
	IssuerID    financing.ID    `gorm:"type:uuid"`
	AskingPrice financing.Money `gorm:"type:float;"`
	WinningBid  *Bid
	Status      InvoiceStatus
}

type Bid struct {
	gorm.Model
	InvoiceID  financing.ID    `gorm:"type:uuid,ForeignKey:InvoiceID"`
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
