package financing

type InvoiceApprovedEvent struct {
	InvoiceID ID `json:"-"`
	Bid       Bid
}

func NewInvoiceApprovedEvent(invoiceID ID, bid Bid) *InvoiceApprovedEvent {
	return &InvoiceApprovedEvent{
		InvoiceID: invoiceID,
		Bid:       bid,
	}
}

func (e *InvoiceApprovedEvent) EventName() string { return "InvoiceApprovedEvent" }

func (e *InvoiceApprovedEvent) WithAggregateID(id string) {
	e.InvoiceID = NewIDFrom(id)
}
