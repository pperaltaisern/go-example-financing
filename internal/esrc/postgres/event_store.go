package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"ledger/internal/esrc"
	"ledger/pkg/financing"

	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v4/pgxpool"
)

type EventStore struct {
	pool *pgxpool.Pool
}

var _ (esrc.EventStore) = (*EventStore)(nil)

func NewEventStore(pool *pgxpool.Pool) *EventStore {
	return &EventStore{
		pool: pool,
	}
}

// event streams

// events
const queryEvents = "SELECT name, data FROM events WHERE event_source_id = $1 ORDER BY version ASC"

func (es *EventStore) Load(ctx context.Context, id esrc.ID) ([]esrc.Event, error) {

}

func (es *EventStore) Contains(ctx context.Context, id esrc.ID) (bool, error) {
	const existsStream = "SELECT EXISTS(SELECT 1 FROM event_streams WHERE id=$1)"
	rows, err := es.pool.Query(ctx, existsStream, id)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	for rows.Next() {
		var exists bool
		err := rows.Scan(&exists)
		if err != nil {
			return false, err
		}
		return exists, nil
	}
	if err = rows.Err(); err != nil {
		return false, err
	}

	return false, nil
}

func (es *EventStore) Create(ctx context.Context, a esrc.Aggregate, t esrc.AggregateType) error {
	return es.pool.BeginFunc(ctx, func(tx pgx.Tx) error {
		const insertEventStream = "INSERT INTO event_streams(id, type, version) VALUES ($1, $2, 0)"
		_, err := tx.Exec(ctx, insertEventStream, a.ID(), t)
		if err != nil {
			return err
		}

		return appendEvents(ctx, tx, a)
	})
}

func (es *EventStore) AppendEvents(ctx context.Context, a esrc.Aggregate) error {
	if a.Version() == len(a.Events()) {
		return errors.New("no new events")
	}

	return es.pool.BeginFunc(ctx, func(tx pgx.Tx) error {
		return appendEvents(ctx, tx, a)
	})
}

func appendEvents(ctx context.Context, tx pgx.Tx, a esrc.Aggregate) error {
	v := a.Version()
	id := a.ID()
	for i, e := range a.Events() {
		const updateStream = "UPDATE event_streams SET version = version + 1 WHERE id = $1 AND VERSION = $2"
		_, err := tx.Exec(ctx, updateStream, id, v)
		if err != nil {
			return err
		}

		b, err := json.Marshal(e)
		if err != nil {
			return err
		}

		const insertEvents = "INSERT INTO events(event_source_id, name, version, data) VALUES ($1, $2, $3, $4)"
		_, err = tx.Exec(ctx, insertEvents, id, e.Name(), v+i+1, b)
		if err != nil {
			return err
		}
	}
	return nil
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
