package financing

type InvoiceFinancedEvent struct {
	InvoiceID   ID
	AskingPrice Money
	Bid         Bid
}

func NewInvoiceFinancedEvent(invoiceID ID, askingPrice Money, bid Bid) InvoiceFinancedEvent {
	return InvoiceFinancedEvent{
		InvoiceID:   invoiceID,
		AskingPrice: askingPrice,
		Bid:         bid,
	}
}

func (e InvoiceFinancedEvent) IsEvent() {}

func (e InvoiceFinancedEvent) Name() string { return "InvoiceFinancedEvent" }
