package financing

import (
	"testing"

	"github.com/pperaltaisern/esrc"

	"github.com/stretchr/testify/require"
)

func TestNewInvoice(t *testing.T) {
	id := NewID()
	issuerID := NewID()
	askingPrice := Money(100)
	inv := NewInvoice(id, issuerID, askingPrice)
	assertNewInvoice(t, inv, id, issuerID, askingPrice)

	e := NewInvoiceCreatedEvent(id, issuerID, askingPrice)
	require.Equal(t, e, inv.Changes()[0])
	require.Equal(t, 0, inv.InitialVersion())
}

func TestNewInvoiceFromEvents(t *testing.T) {
	id := NewID()
	issuerID := NewID()
	askingPrice := Money(100)
	e := NewInvoiceCreatedEvent(id, issuerID, askingPrice)

	inv, err := invoiceFactory{}.NewAggregateFromEvents([]esrc.Event{e})
	require.NoError(t, err)
	assertNewInvoice(t, inv, id, issuerID, askingPrice)

	require.Empty(t, inv.Changes())
	require.Equal(t, 1, inv.InitialVersion())
}

func assertNewInvoice(t *testing.T, inv *Invoice, id, issuerID ID, askingPrice Money) {
	require.Equal(t, id, inv.id)
	require.Equal(t, issuerID, inv.issuerID)
	require.Equal(t, askingPrice, inv.askingPrice)
	require.Equal(t, invoiceStatusAvailable, inv.status)
	require.Nil(t, inv.winningBid)
}

func TestInvoice_ProcessBid(t *testing.T) {
	askingPrice := Money(100)
	tests := map[string]struct {
		bidAmount Money
	}{
		"bid for exactly the asking price": {bidAmount: askingPrice},
		"bid for higher the asking price":  {bidAmount: askingPrice + 1},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			inv := testNewInvoice(100)
			bid := NewBid(NewID(), tc.bidAmount)
			inv.ProcessBid(bid)
			require.Equal(t, invoiceStatusFinanced, inv.status)
			require.Equal(t, bid, *inv.winningBid)

			e := NewInvoiceFinancedEvent(inv.id, inv.askingPrice, bid)
			require.Equal(t, e, inv.Changes()[1])
			require.Equal(t, 0, inv.InitialVersion())
		})
	}
}

func TestInvoice_ProcessBid_NotMatching(t *testing.T) {
	askingPrice := Money(100)
	inv := testNewInvoice(100)

	bid := NewBid(NewID(), askingPrice-1)
	inv.ProcessBid(bid)
	// state of invoice shouldn't have changed
	assertNewInvoice(t, inv, inv.id, inv.issuerID, askingPrice)

	e := NewBidOnInvoiceRejectedEvent(inv.id, bid)
	require.Len(t, inv.Changes(), 2)
	require.Equal(t, e, inv.Changes()[1])
	require.Equal(t, 0, inv.InitialVersion())
}

func TestInvoice_ProcessBid_AfterAlreadyFinanced(t *testing.T) {
	askingPrice := Money(100)
	inv := testNewFinancedInvoice(askingPrice)

	newBid := NewBid(NewID(), askingPrice*2)
	inv.ProcessBid(newBid)

	e := NewBidOnInvoiceRejectedEvent(inv.id, newBid)
	require.Len(t, inv.Changes(), 3)
	require.Equal(t, e, inv.Changes()[2])
	require.Equal(t, 0, inv.InitialVersion())
}

func TestInvoice_ApproveFinancing(t *testing.T) {
	inv := testNewFinancedInvoice(100)
	inv.ApproveFinancing()

	e := NewInvoiceApprovedEvent(inv.id, inv.winningBid.Amount, *inv.winningBid)
	require.Len(t, inv.Changes(), 3)
	require.Equal(t, e, inv.Changes()[2])
	require.Equal(t, 0, inv.InitialVersion())
}

func TestInvoice_ReverseFinancing(t *testing.T) {
	inv := testNewFinancedInvoice(100)
	inv.ReverseFinancing()
	require.Equal(t, invoiceStatusReversed, inv.status)

	e := NewInvoiceReversedEvent(inv.id, inv.askingPrice, *inv.winningBid)
	require.Len(t, inv.Changes(), 3)
	require.Equal(t, e, inv.Changes()[2])
	require.Equal(t, 0, inv.InitialVersion())
}

func TestInvoice_ProcessBid_AfterApproved(t *testing.T) {
	askingPrice := Money(100)
	inv := testNewFinancedInvoice(askingPrice)
	inv.ApproveFinancing()

	bid := NewBid(NewID(), askingPrice)
	inv.ProcessBid(bid)

	e := NewBidOnInvoiceRejectedEvent(inv.id, bid)
	require.Len(t, inv.Changes(), 4)
	require.Equal(t, e, inv.Changes()[3])
	require.Equal(t, 0, inv.InitialVersion())
}

func TestInvoice_ProcessBid_AfterReversed(t *testing.T) {
	askingPrice := Money(100)
	inv := testNewFinancedInvoice(askingPrice)
	inv.ReverseFinancing()

	bid := NewBid(NewID(), askingPrice)
	inv.ProcessBid(bid)

	e := NewBidOnInvoiceRejectedEvent(inv.id, bid)
	require.Len(t, inv.Changes(), 4)
	require.Equal(t, e, inv.Changes()[3])
	require.Equal(t, 0, inv.InitialVersion())
}

func testNewInvoice(askingPrice Money) *Invoice {
	id := NewID()
	issuerID := NewID()
	return NewInvoice(id, issuerID, askingPrice)
}

func testNewFinancedInvoice(askingPrice Money) *Invoice {
	inv := testNewInvoice(askingPrice)
	bid := NewBid(NewID(), askingPrice)
	inv.ProcessBid(bid)
	return inv
}
