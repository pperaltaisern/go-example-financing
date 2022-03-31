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

		return updateInvestor(tx, e.InvestorID, func(investor *Investor) {
			investor.Balance += e.Amount
		}).Error
	})
}

func (c *EventProjector) ProjectBidOnInvoicePlacedEvent(e *financing.BidOnInvoicePlacedEvent) error {
	return c.db.Transaction(func(tx *gorm.DB) error {

		return updateInvestor(tx, e.InvestorID, func(investor *Investor) {
			investor.Balance -= e.BidAmount
			investor.Reserved += e.BidAmount
		}).Error
	})
}

func (c *EventProjector) ProjectInvestorFundsReleasedEvent(e *financing.InvestorFundsReleasedEvent) error {
	return c.db.Transaction(func(tx *gorm.DB) error {

		return updateInvestor(tx, e.InvestorID, func(investor *Investor) {
			investor.Balance += e.Amount
			investor.Reserved -= e.Amount
		}).Error
	})
}

func (c *EventProjector) ProjectInvestorFundsCommittedEvent(e *financing.InvestorFundsCommittedEvent) error {
	return c.db.Transaction(func(tx *gorm.DB) error {

		return updateInvestor(tx, e.InvestorID, func(investor *Investor) {
			investor.Reserved -= e.Amount
			investor.Committed += e.Amount
		}).Error
	})
}

func (c *EventProjector) ProjectInvoiceCreatedEvent(e *financing.InvoiceCreatedEvent) error {
	invoice := &Invoice{
		Invoice: query.Invoice{
			ID:          e.InvoiceID,
			IssuerID:    e.IssuerID,
			AskingPrice: e.AskingPrice,
			Status:      query.InvoiceStatusAvailable,
		},
	}
	tx := c.db.Create(invoice)
	return tx.Error
}

func (c *EventProjector) ProjectInvoiceFinancedEvent(e *financing.InvoiceFinancedEvent) error {
	return c.db.Transaction(func(tx *gorm.DB) error {

		return updateInvoice(tx, e.InvoiceID, func(invoice *Invoice) {
			invoice.Status = query.InvoiceStatusFinanced
			invoice.WinningBid = &query.Bid{
				InvestorID: e.Bid.InvestorID,
				Amount:     e.Bid.Amount,
			}
		}).Error
	})
}
func (c *EventProjector) ProjectInvoiceReversedEvent(e *financing.InvoiceReversedEvent) error {
	return c.db.Transaction(func(tx *gorm.DB) error {

		return updateInvoice(tx, e.InvoiceID, func(invoice *Invoice) {
			invoice.Status = query.InvoiceStatusReversed
		}).Error
	})
}
func (c *EventProjector) ProjectInvoiceApprovedEvent(e *financing.InvoiceApprovedEvent) error {
	return c.db.Transaction(func(tx *gorm.DB) error {

		return updateInvoice(tx, e.InvoiceID, func(invoice *Invoice) {
			invoice.Status = query.InvoiceStatusApproved
		}).Error
	})
}
func (c *EventProjector) ProjectIssuerCreatedEvent(e *financing.IssuerCreatedEvent) error {
	return nil
}

func updateInvestor(tx *gorm.DB, id financing.ID, investorUpdate func(investor *Investor)) *gorm.DB {
	investor := &Investor{}

	tx = tx.
		Clauses(clause.Locking{Strength: "UPDATE"}).
		First(investor, "id = ?", id.String())
	if tx.Error != nil {
		return tx
	}

	investorUpdate(investor)

	return tx.Save(investor)
}

func updateInvoice(tx *gorm.DB, id financing.ID, invoiceUpdate func(invoice *Invoice)) *gorm.DB {
	invoice := &Invoice{}

	tx = tx.
		Clauses(clause.Locking{Strength: "UPDATE"}).
		First(invoice, "id = ?", id.String())
	if tx.Error != nil {
		return tx
	}

	invoiceUpdate(invoice)

	return tx.Save(invoice)
}
