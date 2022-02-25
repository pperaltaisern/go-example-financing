package pg

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/pperaltaisern/financing/pkg/financing"
	"github.com/pperaltaisern/financing/pkg/projection"
	"github.com/pperaltaisern/financing/pkg/query"
)

var _ projection.EventProjector = (*EventProjector)(nil)

type EventProjector struct {
	db *gorm.DB
}

var tables = []interface{}{
	&Investor{},
	&Invoice{},
}

func NewEventProjector(db *gorm.DB) (*EventProjector, error) {
	err := db.AutoMigrate(tables...)
	if err != nil {
		return nil, err
	}
	return &EventProjector{
		db: db,
	}, nil
}

func (c *EventProjector) Clean() error {
	return c.db.Migrator().DropTable(tables...)
}

func (c *EventProjector) ProjectInvestorCreatedEvent(e *financing.InvestorCreatedEvent) error {
	investor := &Investor{
		Investor: query.Investor{
			ID: e.InvestorID,
		},
	}
	tx := c.db.Create(investor)
	return tx.Error
}

func (c *EventProjector) ProjectInvestorFundsAddedEvent(e *financing.InvestorFundsAddedEvent) error {
	return c.db.Transaction(func(tx *gorm.DB) error {
		investor := &Investor{}

		tx = tx.
			Clauses(clause.Locking{Strength: "UPDATE"}).
			First(investor, e.InvestorID)
		if tx.Error != nil {
			return tx.Error
		}

		investor.Balance += e.Amount

		tx = tx.Save(investor)
		return tx.Error
	})
}

func (c *EventProjector) ProjectBidOnInvoicePlacedEvent(e *financing.BidOnInvoicePlacedEvent) error {
	return nil
}

func (c *EventProjector) ProjectInvestorFundsReleasedEvent(e *financing.InvestorFundsReleasedEvent) error {
	return nil
}

func (c *EventProjector) ProjectInvestorFundsCommittedEvent(e *financing.InvestorFundsCommittedEvent) error {
	return nil
}
