package query

import "github.com/pperaltaisern/financing/pkg/financing"

type InvestorQueries interface {
	All() ([]Investor, error)
}

type Investor struct {
	ID        financing.ID    `gorm:"type:uuid"`
	Balance   financing.Money `gorm:"type:float;"`
	Reserved  financing.Money `gorm:"type:float;"`
	Committed financing.Money `gorm:"type:float;"`
}
