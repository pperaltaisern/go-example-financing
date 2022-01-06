package intevent

import (
	"context"
	"ledger/pkg/command"

	"github.com/ThreeDotsLabs/watermill/components/cqrs"
)

type InvestorRegisteredHandler struct {
	bus *cqrs.CommandBus
}

func NewInvestorRegisteredHandler(bus *cqrs.CommandBus) *InvestorRegisteredHandler {
	return &InvestorRegisteredHandler{
		bus: bus,
	}
}

var _ cqrs.EventHandler = (*InvestorRegisteredHandler)(nil)

func (h *InvestorRegisteredHandler) HandlerName() string {
	return "InvestorRegisteredHandler"
}

func (h *InvestorRegisteredHandler) NewEvent() interface{} {
	return &InvestorRegistered{}
}

func (h *InvestorRegisteredHandler) Handle(ctx context.Context, e interface{}) error {
	event := e.(*InvestorRegistered)

	cmd := command.CreateInvestor{
		ID:      event.ID,
		Balance: event.Balance,
	}

	h.bus.Send(ctx, cmd)
	return nil
}
