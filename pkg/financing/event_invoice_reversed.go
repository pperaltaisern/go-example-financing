package financing

type InvoiceReversedEvent struct {
	InvoiceID ID
	Bid       Bid
}

func NewInvoiceReversedEvent(invoiceID ID, bid Bid) InvoiceReversedEvent {
	return InvoiceReversedEvent{
		InvoiceID: invoiceID,
		Bid:       bid,
	}
}

func (e InvoiceReversedEvent) IsEvent() {}

func (e InvoiceReversedEvent) EventName() string { return "InvoiceReversedEvent" }
