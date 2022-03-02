package pg

import (
	"github.com/pperaltaisern/financing/pkg/query"
	"gorm.io/gorm"
)

type Investor struct {
	gorm.Model
	query.Investor
}

type InvestorQueries struct {
	db *gorm.DB
}

var _ query.InvestorQueries = (*InvestorQueries)(nil)

func NewInvestorQueries(db *gorm.DB) *InvestorQueries {
	return &InvestorQueries{
		db: db,
	}
}

func (q *InvestorQueries) All() ([]query.Investor, error) {
	var investors []Investor
	result := q.db.Find(&investors)
	if result.Error != nil {
		return nil, result.Error
	}

	ret := make([]query.Investor, len(investors))
	for i := range investors {
		ret[i] = investors[i].Investor
	}

	return ret, nil
}
