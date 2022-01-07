package financing

type InvoiceFinancedEvent struct {
	InvoiceID   ID `json:"-"`
	AskingPrice Money
	Bid         Bid
}

func NewInvoiceFinancedEvent(invoiceID ID, askingPrice Money, bid Bid) *InvoiceFinancedEvent {
	return &InvoiceFinancedEvent{
		InvoiceID:   invoiceID,
		AskingPrice: askingPrice,
		Bid:         bid,
	}
}

func (e *InvoiceFinancedEvent) EventName() string { return "InvoiceFinancedEvent" }

func (e *InvoiceFinancedEvent) WithAggregateID(id string) {
	e.InvoiceID = NewIDFrom(id)
}
