package financing

type InvoiceApprovedEvent struct {
	InvoiceID ID `json:"-"`
	SoldPrice Money
	Bid       Bid
}

func NewInvoiceApprovedEvent(invoiceID ID, soldPrice Money, bid Bid) *InvoiceApprovedEvent {
	return &InvoiceApprovedEvent{
		InvoiceID: invoiceID,
		SoldPrice: soldPrice,
		Bid:       bid,
	}
}

func (e *InvoiceApprovedEvent) EventName() string { return "InvoiceApprovedEvent" }

func (e *InvoiceApprovedEvent) WithAggregateID(id string) {
	e.InvoiceID = NewIDFromString(id)
}
