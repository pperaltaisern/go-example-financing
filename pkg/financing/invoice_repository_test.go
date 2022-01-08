package financing

import (
	"context"
	"encoding/json"
	"ledger/internal/esrc"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInvoiceRepository_Update(t *testing.T) {
	invoiceID := NewID()
	events := []esrc.Event{
		NewInvoiceCreatedEvent(invoiceID, NewID(), 100),
		NewInvoiceFinancedEvent(invoiceID, 100, NewBid(NewID(), 100)),
	}
	rawEvents := make([]esrc.RawEvent, len(events))
	for i, e := range events {
		b, err := json.Marshal(e)
		require.NoError(t, err)

		rawEvents[i] = esrc.RawEvent{
			Name: e.EventName(),
			Data: b,
		}
	}

	var calls int
	es := &esrc.MockEventStore{
		LoadFn: func(ctx context.Context, id esrc.ID) ([]esrc.RawEvent, error) {
			calls++
			return rawEvents, nil
		},
		AppendEventsFn: func(ctx context.Context, id esrc.ID, fromVersion int, events []esrc.RawEvent) error {
			calls += 10
			require.Equal(t, invoiceID, id)
			require.Equal(t, 2, fromVersion)
			require.Len(t, events, 1)
			return nil
		},
	}
	r := NewInvoiceRepository(es)
	err := r.Update(context.Background(), invoiceID, func(inv *Invoice) error {
		inv.ApproveFinancing()
		return nil
	})

	require.NoError(t, err)
	require.Equal(t, 11, calls)
}
