package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"

	"ledger/internal/esrc"
	"ledger/pkg/financing"
)

const eventStreamInitialVersion = 0

const streamTypeInvestor = "investor"

// event streams
const insertEventStream = "INSERT INTO event_streams(id, type, version) VALUES ($1, $2, $3) ON CONFLICT (id) DO NOTHING"
const existsStream = "SELECT EXISTS(SELECT 1 FROM event_streams WHERE id=$1 AND type=$2)"
const updateStream = "UPDATE event_streams SET version = version + 1 WHERE id = $1 AND VERSION = $2"

// events
const insertEvents = "INSERT INTO events(event_source_id, name, version, data) VALUES ($1, $2, $3, $4) ON CONFLICT (id) DO NOTHING"
const queryEvents = "SELECT name, data FROM events WHERE event_source_id = $1 ORDER BY version ASC"

type InvestorRepository struct {
	pool *pgxpool.Pool
}

var _ financing.InvestorRepository = (*InvestorRepository)(nil)

func NewInvestorRepository(pool *pgxpool.Pool) *InvestorRepository {
	return &InvestorRepository{
		pool: pool,
	}
}

func (r *InvestorRepository) ByID(ctx context.Context, id financing.ID) (*financing.Investor, error) {
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

	var events []esrc.Event
	for rows.Next() {
		var name string
		var data []byte
		err := rows.Scan(&name, &data)
		if err != nil {
			return nil, err
		}

		switch name {
		case "InvestorCreatedEvent":
			var e financing.InvestorCreatedEvent
			err = json.Unmarshal(data, &e)
			if err != nil {
				return nil, err
			}
			events = append(events, e)
		case "InvestorFundsAddedEvent":
			var e financing.InvestorFundsAddedEvent
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
		case "InvestorFundsReleasedEvent":
			var e financing.InvestorFundsReleasedEvent
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

	return financing.NewInvestorFromEvents(events), nil
}

func (r *InvestorRepository) Update(ctx context.Context, inv *financing.Investor) error {
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

func (r *InvestorRepository) Add(ctx context.Context, inv *financing.Investor) error {
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

	_, err = tx.Exec(ctx, insertEventStream, inv.ID(), streamTypeInvestor, eventStreamInitialVersion)
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
