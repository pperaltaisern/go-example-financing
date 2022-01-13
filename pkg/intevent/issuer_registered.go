package intevent

import (
	"github.com/pperaltaisern/financing/pkg/financing"
)

type IssuerRegistered struct {
	ID   financing.ID
	Name string
}
