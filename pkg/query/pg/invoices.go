package pg

import (
	"github.com/pperaltaisern/financing/pkg/query"
	"gorm.io/gorm"
)

type Invoice struct {
	gorm.Model
	query.Invoice
}
