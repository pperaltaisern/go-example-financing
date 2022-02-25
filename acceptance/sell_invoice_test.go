package e2e

import (
	"context"
	"testing"

	"github.com/pperaltaisern/financing/pkg/financing"
	"github.com/pperaltaisern/financing/pkg/grpc/pb"
	"github.com/stretchr/testify/require"
)

func (s *CommandsSuite) TestSellInvoice() {
	t := s.T()

	t.Run("GIVEN an issuer has been registered WHEN he sells an invoice THEN the invoice is created", func(t *testing.T) {
		s.RegisterIssuerAndSellInvoice(20)
	})

	t.Run("GIVEN an unregistered issuer WHEN he sells an invoice THEN the invoice is not created", func(t *testing.T) {
		cmd := &pb.SellInvoiceCommand{
			IssuerId: &pb.UUID{
				Value: financing.NewID().String(),
			},
			AskingPrice: &pb.Money{
				Amount: 20,
			},
		}
		invoiceID, err := s.commands.SellInvoice(context.Background(), cmd)
		require.NoError(t, err)

		s.assertNotContains(t, financing.NewIDFromString(invoiceID.Value))
	})
}

func (s *CommandsSuite) RegisterIssuerAndSellInvoice(askingPrice financing.Money) financing.ID {
	t := s.T()

	// the issuer who sells the invoice must be in database
	issuerID := financing.NewID()
	s.RegisterIssuer(issuerID)
	// send the command
	cmd := &pb.SellInvoiceCommand{
		IssuerId: &pb.UUID{
			Value: issuerID.String(),
		},
		AskingPrice: &pb.Money{
			Amount: float64(askingPrice),
		},
	}
	pbInvoiceID, err := s.commands.SellInvoice(context.Background(), cmd)
	require.NoError(t, err)
	s.assertContains(t, financing.NewIDFromString(pbInvoiceID.Value))

	// assert created events
	invoiceID := financing.NewIDFromString(pbInvoiceID.Value)
	eventAssertion := EventAssertion{
		Expected: financing.NewInvoiceCreatedEvent(invoiceID, issuerID, askingPrice),
		Actual:   &financing.InvoiceCreatedEvent{},
	}
	s.expectEvents(t, eventAssertion)

	return invoiceID
}
