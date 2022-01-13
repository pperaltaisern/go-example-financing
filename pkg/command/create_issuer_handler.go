package command

import (
	"context"

	"github.com/pperaltaisern/financing/pkg/financing"

	"github.com/ThreeDotsLabs/watermill/components/cqrs"
)

type CreateIssuerHandler struct {
	issuers financing.IssuerRepository
}

func NewCreateIssuerHandler(r financing.IssuerRepository) *CreateIssuerHandler {
	return &CreateIssuerHandler{
		issuers: r,
	}
}

var _ cqrs.CommandHandler = (*CreateIssuerHandler)(nil)

func (h *CreateIssuerHandler) HandlerName() string {
	return "CreateIssuerHandler"
}

func (h *CreateIssuerHandler) NewCommand() interface{} {
	return &CreateIssuer{}
}

func (h *CreateIssuerHandler) Handle(ctx context.Context, c interface{}) error {
	cmd := c.(*CreateIssuer)

	issuer := financing.NewIssuer(cmd.ID)

	err := h.issuers.Add(ctx, issuer)
	if err != nil {
		return err
	}

	return nil
}
