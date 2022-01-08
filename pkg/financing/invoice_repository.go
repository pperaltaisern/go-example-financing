package financing

import (
	"context"
	"encoding/json"
	"fmt"
	"ledger/internal/esrc"
)

type InvoiceRepository interface {
	Update(context.Context, ID, UpdateInvoice) error
	Add(context.Context, *Invoice) error
}

type UpdateInvoice func(inv *Invoice) error

func NewInvoiceRepository(es esrc.EventStore) InvoiceRepository {
	return invoiceRepository{es: es}
}

type invoiceRepository struct {
	es esrc.EventStore
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

	newEvents := inv.aggregate.Events()
	if len(newEvents) == 0 {
		return nil
	}
	rawEvents, err := esrc.MarshalEventsJSON(newEvents)
	if err != nil {
		return err
	}
	return r.es.AppendEvents(ctx, inv.id, inv.aggregate.Version(), rawEvents)
}

func (i invoiceRepository) byID(ctx context.Context, id ID) (*Invoice, error) {
	rawEvents, err := i.es.Load(ctx, esrc.ID(id))
	if err != nil {
		return nil, err
	}

	events := make([]esrc.Event, len(rawEvents))
	for i, raw := range rawEvents {
		var e esrc.Event
		switch raw.Name {
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
			return nil, fmt.Errorf("unkown event name: %s", raw.Name)
		}
		err = json.Unmarshal(raw.Data, e)
		if err != nil {
			return nil, err
		}
		events[i] = e
	}

	return newInvoiceFromEvents(events), nil
}

func (i invoiceRepository) Add(ctx context.Context, inv *Invoice) error {
	rawEvents, err := esrc.MarshalEventsJSON(inv.aggregate.Events())
	if err != nil {
		return err
	}

	const aggregateType esrc.AggregateType = "invoice"
	return i.es.Create(ctx, aggregateType, inv.id, rawEvents)
}
