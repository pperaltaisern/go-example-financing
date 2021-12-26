package intevent

import (
	"context"
	"ledger/pkg/command"

	"github.com/ThreeDotsLabs/watermill/components/cqrs"
)

type IssuerRegisteredHandler struct {
	bus *cqrs.CommandBus
}

func NewIssuerRegisteredHandler(bus *cqrs.CommandBus) *IssuerRegisteredHandler {
	return &IssuerRegisteredHandler{
		bus: bus,
	}
}

var _ cqrs.EventHandler = (*IssuerRegisteredHandler)(nil)

func (h *IssuerRegisteredHandler) HandlerName() string {
	return "IssuerRegisteredHandler"
}

func (h *IssuerRegisteredHandler) NewEvent() interface{} {
	return &IssuerRegistered{}
}

func (h *IssuerRegisteredHandler) Handle(ctx context.Context, e interface{}) error {
	event := e.(*IssuerRegistered)

	cmd := command.CreateIssuer{
		ID: event.ID,
	}

	h.bus.Send(ctx, cmd)
	return nil
}
