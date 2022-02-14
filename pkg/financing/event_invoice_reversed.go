package financing

type InvoiceReversedEvent struct {
	InvoiceID ID `json:"-"`
	SoldPrice Money
	Bid       Bid
}

func NewInvoiceReversedEvent(invoiceID ID, soldPrice Money, bid Bid) *InvoiceReversedEvent {
	return &InvoiceReversedEvent{
		InvoiceID: invoiceID,
		SoldPrice: soldPrice,
		Bid:       bid,
	}
}

func (e *InvoiceReversedEvent) EventName() string { return "InvoiceReversedEvent" }

func (e *InvoiceReversedEvent) WithAggregateID(id string) {
	e.InvoiceID = NewIDFromString(id)
}
