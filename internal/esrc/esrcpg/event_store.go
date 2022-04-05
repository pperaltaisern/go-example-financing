package esrcpg

import (
	"context"
	"errors"

	"github.com/pperaltaisern/financing/internal/esrc"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	tableAggregates = "aggregates"
	tableEvents     = "events"
	tableSnapshot   = "snapshots"
)

type EventStore struct {
	pool *pgxpool.Pool
}

var _ esrc.EventStore = (*EventStore)(nil)

func NewEventStore(pool *pgxpool.Pool) *EventStore {
	return &EventStore{
		pool: pool,
	}
}

func (es *EventStore) AddAggregate(ctx context.Context, t esrc.AggregateType, id esrc.ID, events []esrc.RawEvent) error {
	if len(events) == 0 {
		return esrc.ErrAggregateRequiresEvents
	}

	return es.pool.BeginFunc(ctx, func(tx pgx.Tx) error {
		const initialVersion = 0
		const insertEventStream = "INSERT INTO " + tableAggregates + "(id, type, version) VALUES ($1, $2, $3)"
		_, err := tx.Exec(ctx, insertEventStream, id, t, initialVersion)
		if err != nil {
			if isUniqueViolationErr(err) {
				return esrc.ErrAggregateAlreadyExists
			}
		}
		return appendEvents(ctx, tx, id, initialVersion, events)
	})
}

func (es *EventStore) AppendEvents(ctx context.Context, _ esrc.AggregateType, id esrc.ID, fromVersion int, events []esrc.RawEvent) error {
	return es.pool.BeginFunc(ctx, func(tx pgx.Tx) error {
		return appendEvents(ctx, tx, id, fromVersion, events)
	})
}

func appendEvents(ctx context.Context, tx pgx.Tx, id esrc.ID, fromVersion int, events []esrc.RawEvent) error {
	for i, e := range events {
		const updateStream = "UPDATE " + tableAggregates + " SET version = version + 1 WHERE id = $1 AND version = $2"
		cmd, err := tx.Exec(ctx, updateStream, id, fromVersion+i)
		if err != nil {
			return err
		}
		if cmd.RowsAffected() == 0 {
			return esrc.ErrOptimisticConcurrency
		}

		const insertEvents = "INSERT INTO  " + tableEvents + " (aggregate_id, name, version, data) VALUES ($1, $2, $3, $4)"
		_, err = tx.Exec(ctx, insertEvents, id, e.Name, fromVersion+i+1, e.Data)
		if err != nil {
			return err
		}
	}
	return nil
}

func (es *EventStore) AddSnapshot(ctx context.Context, _ esrc.AggregateType, id esrc.ID, snapshot esrc.RawSnapshot) error {
	const selectAggregateVersion = "SELECT version FROM " + tableAggregates + " WHERE id=$1"
	var aggregateVersion int
	err := es.pool.QueryRow(ctx, selectAggregateVersion, id).Scan(&aggregateVersion)
	if err != nil {
		if err == pgx.ErrNoRows {
			return esrc.ErrAggregateNotFound
		}
		return err
	}

	if aggregateVersion < snapshot.Version {
		return esrc.ErrSnapshotWithGreaterVersionThanAggregate
	}

	const insertSnapshot = "INSERT INTO  " + tableSnapshot + " (aggregate_id, version, data) VALUES ($1, $2, $3)"
	_, err = es.pool.Exec(ctx, insertSnapshot, id, snapshot.Version, snapshot.Data)
	return err
}

func (es *EventStore) LatestSnapshot(ctx context.Context, _ esrc.AggregateType, id esrc.ID) (*esrc.RawSnapshot, error) {
	const selectLatestSnapshot = "SELECT version, data FROM " + tableSnapshot + " WHERE aggregate_id=$1 ORDER BY version DESC LIMIT 1"
	var snapshot esrc.RawSnapshot
	err := es.pool.QueryRow(ctx, selectLatestSnapshot, id).Scan(&snapshot.Version, &snapshot.Data)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &snapshot, nil
}

func (es *EventStore) Events(ctx context.Context, _ esrc.AggregateType, id esrc.ID, fromVersion int) ([]esrc.RawEvent, error) {
	const queryEvents = "SELECT name, data FROM  " + tableEvents + "  WHERE aggregate_id = $1 ORDER BY version ASC"
	rows, err := es.pool.Query(ctx, queryEvents, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []esrc.RawEvent
	for rows.Next() {
		var name string
		var data []byte
		err := rows.Scan(&name, &data)
		if err != nil {
			return nil, err
		}
		events = append(events, esrc.RawEvent{Name: name, Data: data})
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	if len(events) == 0 {
		return nil, esrc.ErrAggregateNotFound
	}
	return events, nil
}

func (es *EventStore) ContainsAggregate(ctx context.Context, _ esrc.AggregateType, id esrc.ID) (bool, error) {
	const existsStream = "SELECT EXISTS(SELECT 1 FROM " + tableAggregates + " WHERE id=$1)"
	rows, err := es.pool.Query(ctx, existsStream, id)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	rows.Next()
	err = rows.Err()
	if err != nil {
		return false, err
	}

	var exists bool
	err = rows.Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func isUniqueViolationErr(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation
}
