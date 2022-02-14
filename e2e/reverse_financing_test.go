package e2e

import (
	"context"
	"testing"

	"github.com/pperaltaisern/financing/pkg/financing"
	"github.com/pperaltaisern/financing/pkg/grpc"
	"github.com/pperaltaisern/financing/pkg/grpc/pb"
	"github.com/stretchr/testify/require"
)

func (s *Suite) TestReverseFinancing() {
	t := s.T()

	t.Run(`
	GIVEN an invoice has been financed
	WHEN it's reversed
	THEN an invoice reversed event is raised
	AND the investor's funds reserved for this invoice are released`, func(t *testing.T) {
		result := s.FinanceAnInvoice()

		cmd := &pb.ReverseFinancingCommand{
			InvoiceId: grpc.ConvertID(result.InvoiceID),
		}
		_, err := s.commands.ReverseFinancing(context.Background(), cmd)
		require.NoError(t, err)

		eventAssertions := []EventAssertion{
			{
				Expected: financing.NewInvoiceReversedEvent(result.InvoiceID, result.InvoiceCost, result.Bid),
				Actual:   &financing.InvoiceReversedEvent{},
			},
			{
				Expected: financing.NewInvestorFundsReleasedEvent(result.Bid.InvestorID, result.InvoiceCost),
				Actual:   &financing.InvestorFundsReleasedEvent{},
			},
		}

		s.expectEvents(t, eventAssertions...)
	})
}
