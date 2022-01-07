package eventhandler

import (
	"context"
	"ledger/pkg/financing"

	"github.com/ThreeDotsLabs/watermill/components/cqrs"
)

type InvoiceFinancedHandler struct {
	investors financing.InvestorRepository
}

func NewInvoiceFinancedHandler(r financing.InvestorRepository) *InvoiceFinancedHandler {
	return &InvoiceFinancedHandler{
		investors: r,
	}
}

var _ cqrs.EventHandler = (*InvoiceFinancedHandler)(nil)

func (h *InvoiceFinancedHandler) HandlerName() string {
	return "InvoiceFinancedHandler"
}

func (h *InvoiceFinancedHandler) NewEvent() interface{} {
	return &financing.InvoiceFinancedEvent{}
}

func (h *InvoiceFinancedHandler) Handle(ctx context.Context, e interface{}) error {
	event := e.(*financing.InvoiceFinancedEvent)

	unplacedBid := event.Bid.Amount - event.AskingPrice
	if unplacedBid <= 0 {
		return nil
	}

	return h.investors.Update(ctx, event.Bid.IvestorID, func(investor *financing.Investor) error {
		return investor.ReleaseFunds(unplacedBid)
	})
}
