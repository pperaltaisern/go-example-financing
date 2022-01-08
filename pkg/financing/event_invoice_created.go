package financing

type InvoiceCreatedEvent struct {
	InvoiceID   ID
	IssuerID    ID
	AskingPrice Money
}

func NewInvoiceCreatedEvent(invoiceID ID, issuerID ID, askingPrice Money) *InvoiceCreatedEvent {
	return &InvoiceCreatedEvent{
		InvoiceID:   invoiceID,
		IssuerID:    issuerID,
		AskingPrice: askingPrice,
	}
}

func (e *InvoiceCreatedEvent) EventName() string { return "InvoiceCreatedEvent" }
