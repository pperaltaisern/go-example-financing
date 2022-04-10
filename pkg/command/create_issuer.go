package command

import (
	"context"

	"github.com/pperaltaisern/financing/pkg/financing"
)

type CreateIssuer struct {
	ID financing.ID
}

type CreateIssuerHandler struct {
	issuers financing.IssuerRepository
}

func NewCreateIssuerHandler(r financing.IssuerRepository) *CreateIssuerHandler {
	return &CreateIssuerHandler{
		issuers: r,
	}
}

func (h *CreateIssuerHandler) Handle(ctx context.Context, cmd *CreateIssuer) error {
	issuer := financing.NewIssuer(cmd.ID)
	return h.issuers.Add(ctx, issuer)
}
