package financing

import (
	"testing"

	"github.com/pperaltaisern/financing/internal/esrc"

	"github.com/stretchr/testify/require"
)

func TestNewInvestor(t *testing.T) {
	id := NewID()
	inv := NewInvestor(id)
	assertNewInvestor(t, inv, id)

	e := NewInvestorCreatedEvent(id)
	require.Equal(t, e, inv.Changes()[0])
	require.Equal(t, 0, inv.InitialVersion())
}

func TestNewInvestorFromEvents(t *testing.T) {
	id := NewID()
	e := NewInvestorCreatedEvent(id)

	inv, err := investorFactory{}.NewAggregateFromEvents([]esrc.Event{e})
	require.NoError(t, err)
	assertNewInvestor(t, inv, id)

	require.Empty(t, inv.Changes())
	require.Equal(t, 1, inv.InitialVersion())
}

func assertNewInvestor(t *testing.T, inv *Investor, id ID) {
	require.Equal(t, id, inv.id)
	require.Equal(t, Money(0), inv.balance)
	require.Equal(t, Money(0), inv.reserved)
}

func TestNewInvestorFromSnapshotAndEvents(t *testing.T) {
	id := NewID()
	inv := NewInvestor(id)
	inv.AddFunds(100)
	inv.BidOnInvoice(NewID(), 40)

	snapshot, err := inv.Snapshot()
	require.NoError(t, err)

	rawSnapshot := esrc.RawSnapshot{
		Version: len(inv.Changes()),
		Data:    snapshot,
	}

	invFromSnapshot, err := investorFactory{}.NewAggregateFromSnapshotAndEvents(rawSnapshot, nil)
	require.NoError(t, err)

	require.Equal(t, id, invFromSnapshot.id)
	require.Equal(t, Money(60), invFromSnapshot.balance)
	require.Equal(t, Money(40), invFromSnapshot.reserved)

	require.Empty(t, invFromSnapshot.Changes())
	require.Equal(t, rawSnapshot.Version, invFromSnapshot.InitialVersion())
}

func TestInvestor_AddFunds(t *testing.T) {
	inv := NewInvestor(NewID())
	inv.AddFunds(100)

	require.Equal(t, Money(100), inv.balance)
	require.Equal(t, Money(0), inv.reserved)

	e := NewInvestorFundsAddedEvent(inv.id, 100)
	require.Equal(t, e, inv.Changes()[1])
	require.Equal(t, 0, inv.InitialVersion())
}

func TestInvestor_BidOnInvoice(t *testing.T) {
	balance := Money(100)
	tests := map[string]struct {
		bidAmount Money
	}{
		"entire balance":   {bidAmount: balance},
		"half the balance": {bidAmount: balance / 2},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			inv := testNewInvestorWithBalance(balance)

			invoiceID := NewID()
			err := inv.BidOnInvoice(invoiceID, tc.bidAmount)
			require.NoError(t, err)

			require.Equal(t, balance-tc.bidAmount, inv.balance)
			require.Equal(t, tc.bidAmount, inv.reserved)

			e := NewBidOnInvoicePlacedEvent(inv.id, invoiceID, tc.bidAmount)
			require.Equal(t, e, inv.Changes()[2])
			require.Equal(t, 0, inv.InitialVersion())
		})
	}
}

func TestInvestor_BidOnInvoice_Zero(t *testing.T) {
	balance := Money(100)
	inv := testNewInvestorWithBalance(balance)

	invoiceID := NewID()
	err := inv.BidOnInvoice(invoiceID, 0)
	require.NoError(t, err)

	require.Equal(t, balance, inv.balance)
	require.Equal(t, Money(0), inv.reserved)

	require.Len(t, inv.Changes(), 2)
	require.Equal(t, 0, inv.InitialVersion())
}

func TestInvestor_BidOnInvoice_Negative(t *testing.T) {
	balance := Money(100)
	inv := testNewInvestorWithBalance(balance)

	invoiceID := NewID()
	err := inv.BidOnInvoice(invoiceID, -100)
	require.NoError(t, err)

	require.Equal(t, balance, inv.balance)
	require.Equal(t, Money(0), inv.reserved)

	require.Len(t, inv.Changes(), 2)
	require.Equal(t, 0, inv.InitialVersion())
}

func TestInvestor_BidOnInvoice_HasNotEnoughBalance(t *testing.T) {
	balance := Money(100)
	inv := testNewInvestorWithBalance(balance)

	invoiceID := NewID()
	err := inv.BidOnInvoice(invoiceID, balance+1)
	require.Equal(t, ErrNotEnoughtBalance, err)

	require.Equal(t, balance, inv.balance)
	require.Equal(t, Money(0), inv.reserved)

	require.Len(t, inv.Changes(), 2)
	require.Equal(t, 0, inv.InitialVersion())
}

func testNewInvestorWithBalance(balance Money) *Investor {
	inv := NewInvestor(NewID())
	inv.AddFunds(balance)
	return inv
}

func TestInvestor_ReleaseFunds(t *testing.T) {
	initialFunds := Money(100)

	tests := map[string]struct {
		amountReleased Money
	}{
		"entire reserved money":   {amountReleased: initialFunds},
		"half the reserved money": {amountReleased: initialFunds / 2},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			inv := testNewInvestorWithBalanceAndReserved(t, initialFunds, initialFunds)

			err := inv.ReleaseFunds(tc.amountReleased)
			require.NoError(t, err)

			require.Equal(t, tc.amountReleased, inv.balance)
			require.Equal(t, initialFunds-tc.amountReleased, inv.reserved)

			e := NewInvestorFundsReleasedEvent(inv.id, tc.amountReleased)
			require.Equal(t, e, inv.Changes()[3])
			require.Equal(t, 0, inv.InitialVersion())
		})
	}
}

func TestInvesto_ReleaseFunds_Zero(t *testing.T) {
	initialFunds := Money(100)
	inv := testNewInvestorWithBalanceAndReserved(t, initialFunds, initialFunds)

	err := inv.ReleaseFunds(0)
	require.NoError(t, err)

	require.Equal(t, Money(0), inv.balance)
	require.Equal(t, initialFunds, inv.reserved)

	require.Len(t, inv.Changes(), 3)
	require.Equal(t, 0, inv.InitialVersion())
}

func TestInvestor_ReleaseFunds_Negative(t *testing.T) {
	initialFunds := Money(100)
	inv := testNewInvestorWithBalanceAndReserved(t, initialFunds, initialFunds)

	err := inv.ReleaseFunds(-100)
	require.NoError(t, err)

	require.Equal(t, Money(0), inv.balance)
	require.Equal(t, initialFunds, inv.reserved)

	require.Len(t, inv.Changes(), 3)
	require.Equal(t, 0, inv.InitialVersion())
}

func TestInvestor_ReleaseFunds_HasNotEnoughFundsReserved(t *testing.T) {
	initialFunds := Money(100)
	inv := testNewInvestorWithBalanceAndReserved(t, initialFunds, initialFunds)

	err := inv.ReleaseFunds(initialFunds + 1)
	require.Equal(t, ErrNotEnoughReservedFunds, err)

	require.Equal(t, Money(0), inv.balance)
	require.Equal(t, initialFunds, inv.reserved)

	require.Len(t, inv.Changes(), 3)
	require.Equal(t, 0, inv.InitialVersion())
}

func TestInvestor_CommitFunds(t *testing.T) {
	initialFunds := Money(100)

	tests := map[string]struct {
		amountCommitted Money
	}{
		"entire reserved money":   {amountCommitted: initialFunds},
		"half the reserved money": {amountCommitted: initialFunds / 2},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			inv := testNewInvestorWithBalanceAndReserved(t, initialFunds, initialFunds)

			err := inv.CommitFunds(tc.amountCommitted)
			require.NoError(t, err)

			require.Equal(t, Money(0), inv.balance)
			require.Equal(t, initialFunds-tc.amountCommitted, inv.reserved)

			e := NewInvestorFundsCommittedEvent(inv.id, tc.amountCommitted)
			require.Equal(t, e, inv.Changes()[3])
			require.Equal(t, 0, inv.InitialVersion())
		})
	}
}

func TestInvesto_CommitFunds_Zero(t *testing.T) {
	initialFunds := Money(100)
	inv := testNewInvestorWithBalanceAndReserved(t, initialFunds, initialFunds)

	err := inv.CommitFunds(0)
	require.NoError(t, err)

	require.Equal(t, Money(0), inv.balance)
	require.Equal(t, initialFunds, inv.reserved)

	require.Len(t, inv.Changes(), 3)
	require.Equal(t, 0, inv.InitialVersion())
}

func TestInvestor_CommitFunds_Negative(t *testing.T) {
	initialFunds := Money(100)
	inv := testNewInvestorWithBalanceAndReserved(t, initialFunds, initialFunds)

	err := inv.CommitFunds(-100)
	require.NoError(t, err)

	require.Equal(t, Money(0), inv.balance)
	require.Equal(t, initialFunds, inv.reserved)

	require.Len(t, inv.Changes(), 3)
	require.Equal(t, 0, inv.InitialVersion())
}

func TestInvestor_CommmitFunds_HasNotEnoughFundsReserved(t *testing.T) {
	initialFunds := Money(100)
	inv := testNewInvestorWithBalanceAndReserved(t, initialFunds, initialFunds)

	err := inv.CommitFunds(initialFunds + 1)
	require.Equal(t, ErrNotEnoughReservedFunds, err)

	require.Equal(t, Money(0), inv.balance)
	require.Equal(t, initialFunds, inv.reserved)

	require.Len(t, inv.Changes(), 3)
	require.Equal(t, 0, inv.InitialVersion())
}

func testNewInvestorWithBalanceAndReserved(t *testing.T, balance, reserved Money) *Investor {
	inv := testNewInvestorWithBalance(balance)
	err := inv.BidOnInvoice(NewID(), reserved)
	require.NoError(t, err)
	return inv
}
