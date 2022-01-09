package financing

import (
	"context"
	"fmt"
	"ledger/internal/esrc"
)

type InvoiceRepository interface {
	Update(context.Context, ID, UpdateInvoice) error
	Add(context.Context, *Invoice) error
}

type UpdateInvoice func(inv *Invoice) error

type invoiceRepository struct {
	r esrc.Repository
}

func NewInvoiceRepository(es esrc.EventStore) InvoiceRepository {
	return invoiceRepository{
		r: esrc.NewRepository("invoice", es, invoiceEventsFactory{}, esrc.JSONEventMarshaler{}),
	}
}

type invoiceEventsFactory struct{}

func (invoiceEventsFactory) CreateEmptyEvent(name string) (esrc.Event, error) {
	var e esrc.Event
	switch name {
	case "InvoiceCreatedEvent":
		e = &InvoiceCreatedEvent{}
	case "InvoiceFinancedEvent":
		e = &InvoiceFinancedEvent{}
	case "InvoiceReversedEvent":
		e = &InvoiceReversedEvent{}
	case "BidOnInvoicePlacedEvent":
		e = &BidOnInvoicePlacedEvent{}
	case "InvoiceApprovedEvent":
		e = &InvoiceApprovedEvent{}
	default:
		return nil, fmt.Errorf("unkown event name: %s", name)
	}
	return e, nil
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

	return r.r.Update(ctx, id, inv.aggregate.Version(), inv.aggregate.Events())
}

func (r invoiceRepository) byID(ctx context.Context, id ID) (*Invoice, error) {
	events, err := r.r.ByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return newInvoiceFromEvents(events), nil
}

func (r invoiceRepository) Add(ctx context.Context, inv *Invoice) error {
	return r.r.Add(ctx, inv.id, inv.aggregate.Events())
}
