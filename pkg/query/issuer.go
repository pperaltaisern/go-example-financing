package query

import "github.com/pperaltaisern/financing/pkg/financing"

type IssuerQueries interface {
	All() ([]Issuer, error)
}

type Issuer struct {
	ID financing.ID `gorm:"type:uuid;primary_key;"`
}
