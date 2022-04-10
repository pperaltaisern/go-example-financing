package financing

import (
	"fmt"

	"github.com/pperaltaisern/esrc"
)

type investorEventsFactory struct{}

func (investorEventsFactory) CreateEmptyEvent(name string) (esrc.Event, error) {
	var e esrc.Event
	switch name {
	case "InvestorCreatedEvent":
		e = &InvestorCreatedEvent{}
	case "InvestorFundsAddedEvent":
		e = &InvestorFundsAddedEvent{}
	case "BidOnInvoicePlacedEvent":
		e = &BidOnInvoicePlacedEvent{}
	case "InvestorFundsReleasedEvent":
		e = &InvestorFundsReleasedEvent{}
	case "InvestorFundsCommittedEvent":
		e = &InvestorFundsCommittedEvent{}
	default:
		return nil, fmt.Errorf("unkown event name: %s", name)
	}
	return e, nil
}
