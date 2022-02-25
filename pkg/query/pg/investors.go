package pg

import (
	"github.com/pperaltaisern/financing/pkg/query"
	"gorm.io/gorm"
)

type Investor struct {
	gorm.Model
	query.Investor
}
