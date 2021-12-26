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

	investor, err := h.investors.ByID(ctx, event.Bid.IvestorID)
	if err != nil {
		return err
	}

	investor.ReleaseFunds(unplacedBid)

	h.investors.Update(ctx, investor)
	return nil
}
