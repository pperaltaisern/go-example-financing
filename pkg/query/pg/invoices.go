package pg

import (
	"github.com/pperaltaisern/financing/pkg/financing"
	"github.com/pperaltaisern/financing/pkg/query"
	"gorm.io/gorm"
)

type Invoice struct {
	gorm.Model
	query.Invoice
}

type InvoiceQueries struct {
	db *gorm.DB
}

var _ query.InvoiceQueries = (*InvoiceQueries)(nil)

func NewInvoiceQueries(db *gorm.DB) *InvoiceQueries {
	return &InvoiceQueries{
		db: db,
	}
}

func (q *InvoiceQueries) ByID(id financing.ID) (query.Invoice, error) {
	var invoice Invoice
	result := q.db.First(&invoice, id)
	return invoice.Invoice, result.Error
}

func (q *InvoiceQueries) All() ([]query.Invoice, error) {
	var invoices []Invoice
	result := q.db.Find(&invoices)
	if result.Error != nil {
		return nil, result.Error
	}

	ret := make([]query.Invoice, len(invoices))
	for i := range invoices {
		ret[i] = invoices[i].Invoice
	}
	return ret, nil
}
