package financing

type BidOnInvoiceRejectedEvent struct {
	InvoiceID ID `json:"-"`
	Bid       Bid
}

func NewBidOnInvoiceRejectedEvent(invoiceID ID, bid Bid) *BidOnInvoiceRejectedEvent {
	return &BidOnInvoiceRejectedEvent{
		InvoiceID: invoiceID,
		Bid:       bid,
	}
}

func (e *BidOnInvoiceRejectedEvent) EventName() string { return "BidOnInvoiceRejectedEvent" }

func (e *BidOnInvoiceRejectedEvent) WithAggregateID(id string) {
	e.InvoiceID = NewIDFrom(id)
}
