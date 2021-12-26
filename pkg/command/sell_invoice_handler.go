package command

import (
	"context"
	"fmt"
	"ledger/pkg/financing"
)

type SellInvoiceHandler struct {
	issuers  financing.IssuerRepository
	invoices financing.InvoiceRepository
}

func NewSellInvoiceHandler(issuers financing.IssuerRepository, invoices financing.InvoiceRepository) *SellInvoiceHandler {
	return &SellInvoiceHandler{
		issuers:  issuers,
		invoices: invoices,
	}
}

func (h *SellInvoiceHandler) HandlerName() string {
	return "SellInvoiceHandler"
}

func (h *SellInvoiceHandler) NewCommand() interface{} {
	return &SellInvoice{}
}

func (h *SellInvoiceHandler) Handle(ctx context.Context, c interface{}) error {
	cmd := c.(*SellInvoice)

	found, err := h.issuers.Contains(ctx, cmd.IssuerID)
	if err != nil {
		return err
	}
	if !found {
		return fmt.Errorf("issuer %v not found", cmd.IssuerID)
	}

	invoice := financing.NewInvoice(cmd.IssuerID, cmd.AskingPrice)

	err = h.invoices.Add(ctx, invoice)
	if err != nil {
		return err
	}

	return nil
}
