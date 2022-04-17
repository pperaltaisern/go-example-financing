package pg

import (
	"database/sql"
	"time"

	"github.com/pperaltaisern/financing/pkg/query"
	"gorm.io/gorm"
)

type Issuer struct {
	query.Issuer
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt sql.NullTime
}

type IssuerQueries struct {
	db *gorm.DB
}

var _ query.IssuerQueries = (*IssuerQueries)(nil)

func NewIssuerQueries(db *gorm.DB) *IssuerQueries {
	return &IssuerQueries{
		db: db,
	}
}

func (q *IssuerQueries) All() ([]query.Issuer, error) {
	var issuers []Issuer
	result := q.db.Find(&issuers)
	if result.Error != nil {
		return nil, result.Error
	}

	ret := make([]query.Issuer, len(issuers))
	for i := range issuers {
		ret[i] = issuers[i].Issuer
	}

	return ret, nil
}
