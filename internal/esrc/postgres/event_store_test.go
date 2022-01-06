package postgres

import (
	"context"
	"testing"

	"ledger/internal/esrc/esrctesting"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/require"
)

func TestEventStoreAcceptance(t *testing.T) {
	integrationTest(t)

	const connectionString = "user=postgres password=postgres host=localhost port=5432 dbname=postgres pool_max_conns=10"
	pool, err := pgxpool.Connect(context.Background(), connectionString)
	require.NoError(t, err)

	es := NewEventStore(pool)
	esrctesting.NewEventStoreAcceptance(es).Test(t)
}

func integrationTest(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
}
