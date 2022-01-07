package financing

import (
	"errors"
	"ledger/internal/esrc"
)

type Investor struct {
	aggregate esrc.Aggregate

	id       ID
	balance  Money
	reserved Money
}

func NewInvestor(id ID) *Investor {
	inv := &Investor{}
	inv.aggregate = esrc.NewAggregate(inv.onEvent)

	e := NewInvestorCreatedEvent(id)
	inv.aggregate.Raise(e)
	return inv
}

func newInvestorFromEvents(events []esrc.Event) *Investor {
	inv := &Investor{}
	inv.aggregate = esrc.NewAggregateFromEvents(events, inv.onEvent)
	return inv
}

func (inv *Investor) AddFunds(amount Money) error {
	e := NewInvestorFundsAddedEvent(inv.id, amount)
	inv.aggregate.Raise(e)
	return nil
}

var ErrNotEnoughtBalance = errors.New("there isn't enough balance")

func (inv *Investor) BidOnInvoice(invoiceID ID, amount Money) error {
	if !inv.hasEnoughBalance(amount) {
		return ErrNotEnoughtBalance
	}
	bid := NewBid(inv.id, amount)
	e := NewBidOnInvoicePlacedEvent(invoiceID, bid)
	inv.aggregate.Raise(e)
	return nil
}

var ErrNotEnoughtBalanceReservedToRelease = errors.New("there isn't enough balance reserved to release")

func (inv *Investor) ReleaseFunds(amount Money) error {
	if !inv.hasEnoughBalanceReserved(amount) {
		return ErrNotEnoughtBalanceReservedToRelease
	}
	e := NewInvestorFundsReleasedEvent(inv.id, amount)
	inv.aggregate.Raise(e)
	return nil
}

func (inv *Investor) hasEnoughBalance(amount Money) bool {
	return inv.balance >= amount
}

func (inv *Investor) hasEnoughBalanceReserved(amount Money) bool {
	return inv.reserved < amount
}

func (inv *Investor) addFunds(amount Money) {
	inv.balance += amount
}

func (inv *Investor) reserveFunds(amount Money) {
	inv.balance -= amount
	inv.reserved += amount
}

func (inv *Investor) releaseFunds(amount Money) {
	inv.balance += amount
	inv.reserved -= amount
}

func (inv *Investor) onEvent(event esrc.Event) {
	switch e := event.(type) {
	case *InvestorCreatedEvent:
		inv.id = e.InvestorID
	case *InvestorFundsAddedEvent:
		inv.addFunds(e.Amount)
	case *BidOnInvoicePlacedEvent:
		inv.reserveFunds(e.Bid.Amount)
	case *InvestorFundsReleasedEvent:
		inv.releaseFunds(e.Amount)
	}
}
