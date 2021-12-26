package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"

	"ledger/internal/es"
	"ledger/pkg/financing"
)

const streamTypeInvoice = "Invoice"

type InvoiceRepository struct {
	pool *pgxpool.Pool
}

var _ financing.InvoiceRepository = (*InvoiceRepository)(nil)

func NewInvoiceRepository(pool *pgxpool.Pool) *InvoiceRepository {
	return &InvoiceRepository{
		pool: pool,
	}
}

func (r *InvoiceRepository) ByID(ctx context.Context, id financing.ID) (*financing.Invoice, error) {
	conn, err := r.pool.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	rows, err := conn.Query(ctx, queryEvents, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []es.Event
	for rows.Next() {
		var name string
		var data []byte
		err := rows.Scan(&name, &data)
		if err != nil {
			return nil, err
		}

		switch name {
		case "InvoiceFinancedEvent":
			var e financing.InvoiceFinancedEvent
			err = json.Unmarshal(data, &e)
			if err != nil {
				return nil, err
			}
			events = append(events, e)
		case "InvoiceReversedEvent":
			var e financing.InvoiceReversedEvent
			err = json.Unmarshal(data, &e)
			if err != nil {
				return nil, err
			}
			events = append(events, e)
		case "BidOnInvoicePlacedEvent":
			var e financing.BidOnInvoicePlacedEvent
			err = json.Unmarshal(data, &e)
			if err != nil {
				return nil, err
			}
			events = append(events, e)
		case "InvoiceApprovedEvent":
			var e financing.InvoiceApprovedEvent
			err = json.Unmarshal(data, &e)
			if err != nil {
				return nil, err
			}
			events = append(events, e)
		default:
			return nil, fmt.Errorf("unkown event name: %s", name)
		}
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return financing.NewInvoiceFromEvents(events), nil
}

func (r *InvoiceRepository) Update(ctx context.Context, inv *financing.Invoice) error {
	v := inv.Version()
	events := inv.Events()
	if v == len(events) {
		return errors.New("no new events")
	}

	conn, err := r.pool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	tx, err := conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	for i, e := range events[v:] {
		_, err := tx.Exec(ctx, updateStream, inv.ID(), v)
		if err != nil {
			return err
		}

		b, err := json.Marshal(e)
		if err != nil {
			return err
		}
		_, err = tx.Exec(ctx, insertEvents, inv.ID(), e.Name(), i+1, b)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *InvoiceRepository) Add(ctx context.Context, inv *financing.Invoice) error {
	conn, err := r.pool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	tx, err := conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, insertEventStream, inv.ID(), streamTypeInvoice, eventStreamInitialVersion)
	if err != nil {
		return err
	}

	for i, e := range inv.Events() {
		_, err = tx.Exec(ctx, updateStream, inv.ID(), i)
		if err != nil {
			return err
		}

		b, err := json.Marshal(e)
		if err != nil {
			return err
		}
		_, err = tx.Exec(ctx, insertEvents, inv.ID(), e.Name(), i+1, b)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}
