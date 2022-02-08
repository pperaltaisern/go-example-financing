package e2e

import (
	"context"
	"testing"

	"github.com/pperaltaisern/financing/pkg/financing"
	"github.com/pperaltaisern/financing/pkg/grpc/pb"
	"github.com/stretchr/testify/require"
)

func (s *Suite) TestSellInvoice() {
	t := s.T()

	t.Run("GIVEN a registered issuer WHEN he sells an invoice THEN the invoice is created", func(t *testing.T) {
		// the issuer who sells the invoice must be in database
		issuerID := financing.NewID()
		s.RegisterIssuer(issuerID)
		// send the command
		cmd := &pb.SellInvoiceCommand{
			IssuerId: &pb.UUID{
				Value: issuerID.String(),
			},
			AskingPrice: &pb.Money{
				Amount: 20,
			},
		}
		invoiceID, err := s.commands.SellInvoice(context.Background(), cmd)
		require.NoError(t, err)
		s.assertContains(t, financing.NewIDFromString(invoiceID.Value))

		// assert created events
		eventAssertion := EventAssertion{
			Expected: financing.NewInvoiceCreatedEvent(financing.NewIDFromString(invoiceID.Value), issuerID, 20),
			Actual:   &financing.InvoiceCreatedEvent{},
		}
		s.expectEvents(t, eventAssertion)
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
