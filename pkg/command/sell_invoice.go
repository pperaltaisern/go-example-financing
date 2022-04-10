package command

import (
	"context"
	"fmt"

	"github.com/pperaltaisern/financing/pkg/financing"
)

type SellInvoice struct {
	InvoiceID   financing.ID
	IssuerID    financing.ID
	AskingPrice financing.Money
}

func NewSellInvoice(invoiceID, issuerID financing.ID, askingPrice financing.Money) *SellInvoice {
	return &SellInvoice{
		InvoiceID:   invoiceID,
		IssuerID:    issuerID,
		AskingPrice: askingPrice,
	}
}

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

func (h *SellInvoiceHandler) Handle(ctx context.Context, cmd *SellInvoice) error {
	found, err := h.issuers.Contains(ctx, cmd.IssuerID)
	if err != nil {
		return err
	}
	if !found {
		return fmt.Errorf("issuer %v not found", cmd.IssuerID)
	}

	invoice := financing.NewInvoice(cmd.InvoiceID, cmd.IssuerID, cmd.AskingPrice)
	return h.invoices.Add(ctx, invoice)
}
