package financing

import (
	"errors"

	"github.com/pperaltaisern/financing/internal/esrc"
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

func (inv *Investor) AddFunds(amount Money) {
	e := NewInvestorFundsAddedEvent(inv.id, amount)
	inv.aggregate.Raise(e)
}

var ErrNotEnoughtBalance = errors.New("there isn't enough balance")

func (inv *Investor) BidOnInvoice(invoiceID ID, amount Money) error {
	if amount <= 0 {
		return nil
	}
	if !inv.hasEnoughBalance(amount) {
		return ErrNotEnoughtBalance
	}
	e := NewBidOnInvoicePlacedEvent(inv.id, invoiceID, amount)
	inv.aggregate.Raise(e)
	return nil
}

var ErrNotEnoughReservedFundsToRelease = errors.New("there isn't enough balance reserved to release")

func (inv *Investor) ReleaseFunds(amount Money) error {
	if amount <= 0 {
		return nil
	}
	if !inv.hasEnoughReservedFunds(amount) {
		return ErrNotEnoughReservedFundsToRelease
	}
	e := NewInvestorFundsReleasedEvent(inv.id, amount)
	inv.aggregate.Raise(e)
	return nil
}

func (inv *Investor) hasEnoughBalance(amount Money) bool {
	return inv.balance >= amount
}

func (inv *Investor) hasEnoughReservedFunds(amount Money) bool {
	return inv.reserved >= amount
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
		inv.reserveFunds(e.BidAmount)
	case *InvestorFundsReleasedEvent:
		inv.releaseFunds(e.Amount)
	}
}
