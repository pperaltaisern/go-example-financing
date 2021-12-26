package financing

type BidOnInvoiceRejectedEvent struct {
	InvoiceID ID
	Bid       Bid
}

func NewBidOnInvoiceRejectedEvent(invoiceID ID, bid Bid) BidOnInvoiceRejectedEvent {
	return BidOnInvoiceRejectedEvent{
		InvoiceID: invoiceID,
		Bid:       bid,
	}
}

func (e BidOnInvoiceRejectedEvent) IsEvent() {}

func (e BidOnInvoiceRejectedEvent) Name() string { return "BidOnInvoiceRejectedEvent" }
