package eventhandler

import (
	"context"

	"github.com/pperaltaisern/financing/pkg/financing"

	"github.com/ThreeDotsLabs/watermill/components/cqrs"
)

type InvoiceReversedHandler struct {
	investors financing.InvestorRepository
}

func NewInvoiceReversedHandler(r financing.InvestorRepository) *InvoiceReversedHandler {
	return &InvoiceReversedHandler{
		investors: r,
	}
}

var _ cqrs.EventHandler = (*InvoiceReversedHandler)(nil)

func (h *InvoiceReversedHandler) HandlerName() string {
	return "InvoiceReversedHandler"
}

func (h *InvoiceReversedHandler) NewEvent() interface{} {
	return &financing.InvoiceReversedEvent{}
}

func (h *InvoiceReversedHandler) Handle(ctx context.Context, e interface{}) error {
	event := e.(*financing.InvoiceReversedEvent)

	return h.investors.Update(ctx, event.Bid.IvestorID, func(investor *financing.Investor) error {
		return investor.ReleaseFunds(event.Bid.Amount)
	})
}
