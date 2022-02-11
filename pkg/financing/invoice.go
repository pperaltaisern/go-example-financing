package financing

import (
	"github.com/pperaltaisern/financing/internal/esrc"
)

type Invoice struct {
	aggregate esrc.Aggregate

	id          ID
	issuerID    ID
	askingPrice Money
	status      invoiceStatus
	winningBid  *Bid
}

type invoiceStatus byte

const (
	invoiceStatusAvailable invoiceStatus = iota
	invoiceStatusFinanced
	invoiceStatusApproved
	invoiceStatusReversed
)

func NewInvoice(id, issuerID ID, askingPrice Money) *Invoice {
	inv := &Invoice{}
	inv.aggregate = esrc.NewAggregate(inv.onEvent)

	e := NewInvoiceCreatedEvent(id, issuerID, askingPrice)
	inv.aggregate.Raise(e)
	return inv
}

func newInvoiceFromEvents(events []esrc.Event) *Invoice {
	inv := &Invoice{}
	inv.aggregate = esrc.NewAggregateFromEvents(events, inv.onEvent)
	return inv
}

func (inv *Invoice) ProcessBid(bid Bid) {
	if inv.status != invoiceStatusAvailable || !inv.isMatchingBid(bid) {
		e := NewBidOnInvoiceRejectedEvent(inv.id, bid)
		inv.aggregate.Raise(e)
		return
	}

	e := NewInvoiceFinancedEvent(inv.id, inv.askingPrice, bid)
	inv.aggregate.Raise(e)
}

func (inv *Invoice) isMatchingBid(bid Bid) bool {
	return bid.Amount >= inv.askingPrice
}

func (inv *Invoice) finance(bid Bid) {
	inv.status = invoiceStatusFinanced
	inv.winningBid = &bid
}

func (inv *Invoice) reverse() {
	inv.status = invoiceStatusReversed
}

func (inv *Invoice) approve() {
	inv.status = invoiceStatusApproved
}

func (inv *Invoice) ReverseFinancing() {
	if inv.status != invoiceStatusFinanced {
		return
	}

	e := NewInvoiceReversedEvent(inv.id, *inv.winningBid)
	inv.aggregate.Raise(e)
}

func (inv *Invoice) ApproveFinancing() {
	if inv.status != invoiceStatusFinanced {
		return
	}

	e := NewInvoiceApprovedEvent(inv.id, inv.askingPrice, *inv.winningBid)
	inv.aggregate.Raise(e)
}

func (inv *Invoice) onEvent(event esrc.Event) {
	switch e := event.(type) {
	case *InvoiceCreatedEvent:
		inv.id = e.InvoiceID
		inv.issuerID = e.IssuerID
		inv.askingPrice = e.AskingPrice
		inv.status = invoiceStatusAvailable
	case *InvoiceFinancedEvent:
		inv.finance(e.Bid)
	case *InvoiceReversedEvent:
		inv.reverse()
	case *InvoiceApprovedEvent:
		inv.approve()
	}
}
