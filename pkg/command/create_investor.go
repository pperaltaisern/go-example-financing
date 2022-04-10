package command

import (
	"context"

	"github.com/pperaltaisern/financing/pkg/financing"
)

type CreateInvestor struct {
	ID      financing.ID
	Balance financing.Money
}

type CreateInvestorHandler struct {
	investors financing.InvestorRepository
}

func NewCreateInvestorHandler(r financing.InvestorRepository) *CreateInvestorHandler {
	return &CreateInvestorHandler{
		investors: r,
	}
}

func (h *CreateInvestorHandler) Handle(ctx context.Context, cmd *CreateInvestor) error {
	investor := financing.NewInvestor(cmd.ID)
	investor.AddFunds(cmd.Balance)

	err := h.investors.Add(ctx, investor)
	if err != nil {
		return err
	}

	return nil
}
