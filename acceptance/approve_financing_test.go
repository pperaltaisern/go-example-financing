package e2e

import (
	"context"
	"testing"

	"github.com/pperaltaisern/financing/pkg/financing"
	"github.com/pperaltaisern/financing/pkg/grpc"
	"github.com/pperaltaisern/financing/pkg/grpc/pb"
	"github.com/stretchr/testify/require"
)

func (s *CommandsSuite) TestApproveFinancing() {
	t := s.T()

	t.Run(`
	GIVEN an invoice has been financed
	WHEN it's approved
	THEN an invoice approved event is raised
	AND the investor's funds reserved for this invoice are commited`, func(t *testing.T) {
		result := s.FinanceAnInvoice()

		cmd := &pb.ApproveFinancingCommand{
			InvoiceId: grpc.ConvertID(result.InvoiceID),
		}
		_, err := s.commands.ApproveFinancing(context.Background(), cmd)
		require.NoError(t, err)

		eventAssertions := []EventAssertion{
			{
				Expected: financing.NewInvoiceApprovedEvent(result.InvoiceID, result.InvoiceCost, result.Bid),
				Actual:   &financing.InvoiceApprovedEvent{},
			},
			{
				Expected: financing.NewInvestorFundsCommittedEvent(result.Bid.InvestorID, result.InvoiceCost),
				Actual:   &financing.InvestorFundsCommittedEvent{},
			},
		}

		s.expectEvents(t, eventAssertions...)
	})
}
