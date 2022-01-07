package financing

type BidOnInvoicePlacedEvent struct {
	InvoiceID ID
	Bid       Bid
}

func NewBidOnInvoicePlacedEvent(invoiceID ID, bid Bid) *BidOnInvoicePlacedEvent {
	return &BidOnInvoicePlacedEvent{
		InvoiceID: invoiceID,
		Bid:       bid,
	}
}

func (e *BidOnInvoicePlacedEvent) EventName() string { return "BidOnInvoicePlacedEvent" }

func (e *BidOnInvoicePlacedEvent) WithAggregateID(id string) {
	e.InvoiceID = NewIDFrom(id)
}
