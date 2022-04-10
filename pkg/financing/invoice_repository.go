package financing

import (
	"context"

	"github.com/pperaltaisern/financing/internal/esrc"
)

type InvoiceRepository interface {
	Update(context.Context, ID, UpdateInvoice) error
	Add(context.Context, *Invoice) error
}

type UpdateInvoice func(inv *Invoice) error

type invoiceRepository struct {
	r *esrc.Repository[*Invoice]
}

func NewInvoiceRepository(es esrc.EventStore, opts ...esrc.RepositoryOption[*Invoice]) InvoiceRepository {
	return invoiceRepository{
		r: esrc.NewRepository[*Invoice](es, invoiceFactory{}, invoiceEventsFactory{}, opts...),
	}
}

func (r invoiceRepository) Update(ctx context.Context, id ID, update UpdateInvoice) error {
	inv, err := r.byID(ctx, id)
	if err != nil {
		return err
	}

	err = update(inv)
	if err != nil {
		return err
	}

	return r.r.Update(ctx, inv)
}

func (r invoiceRepository) byID(ctx context.Context, id ID) (*Invoice, error) {
	return r.r.FindByID(ctx, id)
}

func (r invoiceRepository) Add(ctx context.Context, inv *Invoice) error {
	return r.r.Add(ctx, inv)
}
