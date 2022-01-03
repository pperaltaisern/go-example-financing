package command

import (
	"context"
	"fmt"
	"ledger/pkg/financing"

	"github.com/ThreeDotsLabs/watermill/components/cqrs"
)

type BidOnInvoiceHandler struct {
	investors financing.InvestorRepository
}

func NewBidOnInvoiceHandler(r financing.InvestorRepository) *BidOnInvoiceHandler {
	return &BidOnInvoiceHandler{
		investors: r,
	}
}

var _ cqrs.CommandHandler = (*BidOnInvoiceHandler)(nil)

func (h *BidOnInvoiceHandler) HandlerName() string {
	return "BidOnInvoiceHandler"
}

func (h *BidOnInvoiceHandler) NewCommand() interface{} {
	return &BidOnInvoice{}
}

func (h *BidOnInvoiceHandler) Handle(ctx context.Context, c interface{}) error {
	cmd := c.(*BidOnInvoice)

	fmt.Println("waaaaaaaaaaaat")
	investor, err := h.investors.ByID(ctx, cmd.InvestorID)
	if err != nil {
		return nil
	}

	fmt.Println("investor: ", investor)

	err = investor.BidOnInvoice(cmd.InvoiceID, cmd.BidAmount)
	if err != nil {

		fmt.Println("BidOnInvoice: ", err)
		return err
	}

	fmt.Println("UPdate: ", investor.Version(), len(investor.Events()))
	err = h.investors.Update(ctx, investor)
	if err != nil {
		fmt.Println("UPdate: err", err)
		return err
	}

	return nil
}
