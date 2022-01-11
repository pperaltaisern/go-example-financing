package financing

type BidOnInvoicePlacedEvent struct {
	InvestorID ID `json:"-"`
	InvoiceID  ID
	BidAmount  Money
}

func NewBidOnInvoicePlacedEvent(investorID, invoiceID ID, bidAmount Money) *BidOnInvoicePlacedEvent {
	return &BidOnInvoicePlacedEvent{
		InvestorID: investorID,
		InvoiceID:  invoiceID,
		BidAmount:  bidAmount,
	}
}

func (e *BidOnInvoicePlacedEvent) EventName() string { return "BidOnInvoicePlacedEvent" }

func (e *BidOnInvoicePlacedEvent) WithAggregateID(id string) {
	e.InvestorID = NewIDFromString(id)
}
