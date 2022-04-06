package esrcpg

import (
	"context"
	"testing"

	"github.com/pperaltaisern/financing/internal/esrc/esrctesting"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/require"
)

func TestEventStoreAcceptance(t *testing.T) {
	integrationTest(t)

	const connectionString = "user=postgres password=postgres host=localhost port=5432 dbname=postgres pool_max_conns=10"
	pool, err := pgxpool.Connect(context.Background(), connectionString)
	require.NoError(t, err)

	es := NewEventStore(pool)
	o := NewEventStoreOutbox(pool)
	acceptance := esrctesting.NewEventStoreAcceptanceSuite(es,
		esrctesting.EventStoreAcceptanceSuiteWithOutbox(o))
	acceptance.Test(t)
}

func integrationTest(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
}
