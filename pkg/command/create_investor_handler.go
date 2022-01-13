package command

import (
	"context"
	"fmt"

	"github.com/pperaltaisern/financing/pkg/financing"

	"github.com/ThreeDotsLabs/watermill/components/cqrs"
)

type CreateInvestorHandler struct {
	investors financing.InvestorRepository
}

func NewCreateInvestorHandler(r financing.InvestorRepository) *CreateInvestorHandler {
	return &CreateInvestorHandler{
		investors: r,
	}
}

var _ cqrs.CommandHandler = (*CreateInvestorHandler)(nil)

func (h *CreateInvestorHandler) HandlerName() string {
	return "CreateInvestorHandler"
}

func (h *CreateInvestorHandler) NewCommand() interface{} {
	return &CreateInvestor{}
}

func (h *CreateInvestorHandler) Handle(ctx context.Context, c interface{}) error {
	fmt.Println("HANDLER")
	cmd := c.(*CreateInvestor)

	investor := financing.NewInvestor(cmd.ID)
	investor.AddFunds(cmd.Balance)

	err := h.investors.Add(ctx, investor)
	if err != nil {
		fmt.Println("HANDLER ERR: ", err)
		return err
	}

	return nil
}
