package financing

import (
	"errors"
	"ledger/internal/esrc"
)

type Invoice struct {
	esrc.Aggregate

	id          ID
	issuerID    ID
	askingPrice Money
	status      invoiceStatus
	WinningBid  *Bid
}

type invoiceStatus byte

const (
	invoiceStatusAvailable invoiceStatus = iota
	invoiceStatusFinanced
	invoiceStatusApproved
	invoiceStatusReversed
)

func NewInvoice(id, issuerID ID, askingPrice Money) *Invoice {
	inv := &Invoice{
		id:          id,
		issuerID:    issuerID,
		askingPrice: askingPrice,
		status:      invoiceStatusAvailable,
	}
	inv.Aggregate = esrc.NewAggregate(inv.onEvent)
	return inv
}

func (inv *Invoice) ID() ID {
	return inv.id
}

var ErrBidAmountIsLowerThanTheAskingPrice = errors.New("bid amount is lower than the invoice's asking price")

func (inv *Invoice) ProcessBid(bid Bid) {
	if inv.status != invoiceStatusAvailable || !inv.isMatchingBid(bid) {
		e := NewBidOnInvoiceRejectedEvent(inv.id, bid)
		inv.Raise(e)
		return
	}

	e := NewInvoiceFinancedEvent(inv.id, inv.askingPrice, bid)
	inv.Raise(e)
}

func (inv *Invoice) isMatchingBid(bid Bid) bool {
	return bid.Amount >= inv.askingPrice
}

func (inv *Invoice) finance(bid Bid) {
	inv.status = invoiceStatusFinanced
	inv.WinningBid = &bid
}

func (inv *Invoice) reverse() {
	inv.status = invoiceStatusReversed
	inv.WinningBid = nil
}

func (inv *Invoice) approve() {
	inv.status = invoiceStatusApproved
}

func (inv *Invoice) ReverseFinancing() {
	if inv.status != invoiceStatusFinanced {
		return
	}

	e := NewInvoiceReversedEvent(inv.id, *inv.WinningBid)
	inv.Raise(e)
}

func (inv *Invoice) ApproveFinancing() {
	if inv.status != invoiceStatusFinanced {
		return
	}

	e := NewInvoiceApprovedEvent(inv.id, *inv.WinningBid)
	inv.Raise(e)
}

func NewInvoiceFromEvents(events []esrc.Event) *Invoice {
	inv := &Invoice{}
	inv.Aggregate = esrc.NewAggregate(inv.onEvent)

	inv.Replay(events)

	return inv
}

func (inv *Invoice) onEvent(event esrc.Event) {
	switch e := event.(type) {
	case InvoiceFinancedEvent:
		inv.finance(e.Bid)
	case InvoiceReversedEvent:
		inv.reverse()
	case InvoiceApprovedEvent:
		inv.approve()
	}
}
