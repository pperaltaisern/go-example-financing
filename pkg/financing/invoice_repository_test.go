package financing

import (
	"context"
	"testing"

	"github.com/pperaltaisern/financing/internal/esrc"

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
		LoadFn: func(ctx context.Context, t esrc.AggregateType, id esrc.ID) ([]esrc.RawEvent, error) {
			calls++
			return rawEvents, nil
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
	require.Equal(t, 11, calls)
}
