package intevent

import (
	"context"

	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/pperaltaisern/financing/pkg/command"
	"github.com/pperaltaisern/financing/pkg/financing"
)

type IssuerRegistered struct {
	ID   financing.ID
	Name string
}

type IssuerRegisteredHandler struct {
	bus *cqrs.CommandBus
}

func NewIssuerRegisteredHandler(bus *cqrs.CommandBus) *IssuerRegisteredHandler {
	return &IssuerRegisteredHandler{
		bus: bus,
	}
}

func (h *IssuerRegisteredHandler) Handle(ctx context.Context, event *IssuerRegistered) error {
	cmd := command.CreateIssuer{
		ID: event.ID,
	}

	h.bus.Send(ctx, cmd)
	return nil
}
