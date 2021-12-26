package financing

type InvoiceApprovedEvent struct {
	InvoiceID ID
	Bid       Bid
}

func NewInvoiceApprovedEvent(invoiceID ID, bid Bid) InvoiceApprovedEvent {
	return InvoiceApprovedEvent{
		InvoiceID: invoiceID,
		Bid:       bid,
	}
}

func (e InvoiceApprovedEvent) IsEvent() {}

func (e InvoiceApprovedEvent) Name() string { return "InvoiceApprovedEvent" }
