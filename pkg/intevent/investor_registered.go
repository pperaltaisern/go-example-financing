package intevent

import (
	"context"

	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/pperaltaisern/financing/pkg/command"
	"github.com/pperaltaisern/financing/pkg/financing"
)

type InvestorRegistered struct {
	ID      financing.ID
	Name    string
	Balance financing.Money
}

type InvestorRegisteredHandler struct {
	bus *cqrs.CommandBus
}

func NewInvestorRegisteredHandler(bus *cqrs.CommandBus) *InvestorRegisteredHandler {
	return &InvestorRegisteredHandler{
		bus: bus,
	}
}

func (h *InvestorRegisteredHandler) Handle(ctx context.Context, event *InvestorRegistered) error {
	cmd := command.CreateInvestor{
		ID:      event.ID,
		Balance: event.Balance,
	}

	h.bus.Send(ctx, cmd)
	return nil
}
