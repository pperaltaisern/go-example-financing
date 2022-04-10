package financing

import (
	"encoding/json"
	"errors"

	"github.com/pperaltaisern/financing/internal/esrc"
)

type Investor struct {
	esrc.EventRaiserAggregate

	id       ID
	balance  Money
	reserved Money
}

var _ esrc.Aggregate = (*Investor)(nil)

func NewInvestor(id ID) *Investor {
	inv := &Investor{}
	inv.EventRaiserAggregate = esrc.NewEventRaiserAggregate(inv.onEvent)

	e := NewInvestorCreatedEvent(id)
	inv.Raise(e)
	return inv
}

func (inv *Investor) ID() esrc.ID {
	return inv.id
}

func (inv *Investor) AddFunds(amount Money) {
	e := NewInvestorFundsAddedEvent(inv.id, amount)
	inv.Raise(e)
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
	inv.Raise(e)
	return nil
}

var ErrNotEnoughReservedFunds = errors.New("there isn't enough balance reserved")

func (inv *Investor) ReleaseFunds(amount Money) error {
	if amount <= 0 {
		return nil
	}
	if !inv.hasEnoughReservedFunds(amount) {
		return ErrNotEnoughReservedFunds
	}
	e := NewInvestorFundsReleasedEvent(inv.id, amount)
	inv.Raise(e)
	return nil
}

func (inv *Investor) CommitFunds(amount Money) error {
	if amount <= 0 {
		return nil
	}
	if !inv.hasEnoughReservedFunds(amount) {
		return ErrNotEnoughReservedFunds
	}
	e := NewInvestorFundsCommittedEvent(inv.id, amount)
	inv.Raise(e)
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

func (inv *Investor) commitFunds(amount Money) {
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
	case *InvestorFundsCommittedEvent:
		inv.commitFunds(e.Amount)
	}
}

func (inv *Investor) Snapshot() ([]byte, error) {
	return json.Marshal(investorSnapshot{
		ID:       inv.id,
		Balance:  inv.balance,
		Reserved: inv.reserved,
	})
}

type investorSnapshot struct {
	ID       ID
	Balance  Money
	Reserved Money
}
