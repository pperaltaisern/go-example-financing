package financing

type InvoiceReversedEvent struct {
	InvoiceID ID `json:"-"`
	Bid       Bid
}

func NewInvoiceReversedEvent(invoiceID ID, bid Bid) *InvoiceReversedEvent {
	return &InvoiceReversedEvent{
		InvoiceID: invoiceID,
		Bid:       bid,
	}
}

func (e InvoiceReversedEvent) EventName() string { return "InvoiceReversedEvent" }

func (e InvoiceReversedEvent) WithAggregateID(id string) {
	e.InvoiceID = NewIDFrom(id)
}
