package financing

import (
	"fmt"

	"github.com/pperaltaisern/financing/internal/esrc"
)

type invoiceEventsFactory struct{}

func (invoiceEventsFactory) CreateEmptyEvent(name string) (esrc.Event, error) {
	var e esrc.Event
	switch name {
	case "InvoiceCreatedEvent":
		e = &InvoiceCreatedEvent{}
	case "InvoiceFinancedEvent":
		e = &InvoiceFinancedEvent{}
	case "InvoiceReversedEvent":
		e = &InvoiceReversedEvent{}
	case "BidOnInvoicePlacedEvent":
		e = &BidOnInvoicePlacedEvent{}
	case "InvoiceApprovedEvent":
		e = &InvoiceApprovedEvent{}
	default:
		return nil, fmt.Errorf("unkown event name: %s", name)
	}
	return e, nil
}
