package postgres

import (
	"context"
	"ledger/internal/esrc/relay"

	"github.com/jackc/pgx/v4/pgxpool"
)

type EventStoreOutbox struct {
	pool *pgxpool.Pool
}

var _ (relay.EventStoreOutbox) = (*EventStoreOutbox)(nil)

func NewEventStoreOutbox(pool *pgxpool.Pool) *EventStoreOutbox {
	return &EventStoreOutbox{
		pool: pool,
	}
}

func (o *EventStoreOutbox) UnpublishedEvents(ctx context.Context) ([]relay.Event, error) {
	const query = "SELECT id, name, data FROM events WHERE published = FALSE GROUP BY (event_source_id, id) ORDER BY version ASC"

	conn, err := o.pool.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	rows, err := conn.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []relay.Event
	for rows.Next() {
		var id uint64
		var name string
		var data []byte
		err := rows.Scan(&id, &name, &data)
		if err != nil {
			return nil, err
		}
		re := relay.NewEvent(id, name, data)
		events = append(events, re)
	}
	return events, nil
}

func (o *EventStoreOutbox) MarkEventsAsPublised(ctx context.Context, events []relay.Event) error {
	const update = "UPDATE events SET published = TRUE WHERE id = $1"

	conn, err := o.pool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	tx, err := conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	for _, e := range events {
		_, err = tx.Exec(ctx, update, e.ID)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}
