package eventhandler

import (
	"context"

	"github.com/pperaltaisern/financing/pkg/financing"

	"github.com/ThreeDotsLabs/watermill/components/cqrs"
)

type InvoiceApprovedHandler struct {
	investors financing.InvestorRepository
}

func NewInvoiceApprovedHandler(r financing.InvestorRepository) *InvoiceApprovedHandler {
	return &InvoiceApprovedHandler{
		investors: r,
	}
}

var _ cqrs.EventHandler = (*InvoiceApprovedHandler)(nil)

func (h *InvoiceApprovedHandler) HandlerName() string {
	return "InvoiceApprovedHandler"
}

func (h *InvoiceApprovedHandler) NewEvent() interface{} {
	return &financing.InvoiceApprovedEvent{}
}

func (h *InvoiceApprovedHandler) Handle(ctx context.Context, e interface{}) error {
	event := e.(*financing.InvoiceApprovedEvent)

	return h.investors.Update(ctx, event.Bid.InvestorID, func(investor *financing.Investor) error {
		return investor.CommitFunds(event.SoldPrice)
	})
}
