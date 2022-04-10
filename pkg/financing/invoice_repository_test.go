package financing

import (
	"context"
	"testing"

	"github.com/pperaltaisern/esrc"

	"github.com/stretchr/testify/require"
)

func TestInvoiceRepository_Update(t *testing.T) {
	invoiceID := NewID()
	events := []esrc.Event{
		NewInvoiceCreatedEvent(invoiceID, NewID(), 100),
		NewInvoiceFinancedEvent(invoiceID, 100, NewBid(NewID(), 100)),
	}
	rawEvents, err := esrc.MarshalEvents(events, esrc.JSONEventMarshaler{})
	require.NoError(t, err)

	var calls int
	es := &esrc.MockEventStore{
		EventsFn: func(ctx context.Context, at esrc.AggregateType, i1 esrc.ID, i2 int) ([]esrc.RawEvent, error) {
			calls++
			return rawEvents, nil
		},
		LatestSnapshotFn: func(ctx context.Context, at esrc.AggregateType, i esrc.ID) (*esrc.RawSnapshot, error) {
			calls += 2
			return nil, nil
		},
		AppendEventsFn: func(ctx context.Context, _ esrc.AggregateType, id esrc.ID, fromVersion int, events []esrc.RawEvent) error {
			calls += 10
			require.Equal(t, invoiceID, id)
			require.Equal(t, 2, fromVersion)
			require.Len(t, events, 1)
			return nil
		},
	}

	r := NewInvoiceRepository(es)
	err = r.Update(context.Background(), invoiceID, func(inv *Invoice) error {
		inv.ApproveFinancing()
		return nil
	})

	require.NoError(t, err)
	require.Equal(t, 13, calls)
}
