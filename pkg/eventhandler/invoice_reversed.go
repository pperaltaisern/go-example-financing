package eventhandler

import (
	"context"
	"ledger/pkg/financing"

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

	investor, err := h.investors.ByID(ctx, event.Bid.IvestorID)
	if err != nil {
		return err
	}

	investor.ReleaseFunds(event.Bid.Amount)

	h.investors.Update(ctx, investor)
	return nil
}
