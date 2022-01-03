package postgres

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/jackc/pgx/v4/pgxpool"

	"ledger/pkg/financing"
)

const streamTypeIssuer = "Issuer"

type IssuerRepository struct {
	pool *pgxpool.Pool
}

var _ financing.IssuerRepository = (*IssuerRepository)(nil)

func NewIssuerRepository(pool *pgxpool.Pool) *IssuerRepository {
	return &IssuerRepository{
		pool: pool,
	}
}

func (r *IssuerRepository) Contains(ctx context.Context, id financing.ID) (bool, error) {
	conn, err := r.pool.Acquire(ctx)
	if err != nil {
		return false, err
	}
	defer conn.Release()

	rows, err := conn.Query(ctx, existsStream, id)
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

func (r *IssuerRepository) Update(ctx context.Context, inv *financing.Issuer) error {
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

	for i, e := range events {
		_, err := tx.Exec(ctx, updateStream, inv.ID(), v)
		if err != nil {
			return err
		}

		b, err := json.Marshal(e)
		if err != nil {
			return err
		}
		_, err = tx.Exec(ctx, insertEvents, inv.ID(), e.Name(), v+i+1, b)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *IssuerRepository) Add(ctx context.Context, inv *financing.Issuer) error {
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

	_, err = tx.Exec(ctx, insertEventStream, inv.ID(), streamTypeIssuer, eventStreamInitialVersion)
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
