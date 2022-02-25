package e2e

import (
	"context"
	"testing"

	"github.com/pperaltaisern/financing/pkg/financing"
	"github.com/pperaltaisern/financing/pkg/grpc"
	"github.com/pperaltaisern/financing/pkg/grpc/pb"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *QueriesSuite) TestAllInvoices() {
	t := s.T()

	t.Run(`
	GIVEN that there isn't any invoice
	WHEN all invoices are queried
	THEN no results are obtained`, func(t *testing.T) {
		expected := []pb.Invoice{}

		invoices, err := s.queries.AllInvoices(context.Background(), &emptypb.Empty{})
		require.NoError(t, err)

		require.Equal(t, expected, invoices)
	})

	t.Run(`
	GIVEN there is an issuer registered 
	AND she sells an invoice with asking price of 15
	AND she sells an invoice with asking price of 20
	WHEN all invoices are queried
	THEN 2 available invoices are obtained without bids`, func(t *testing.T) {
		issuerID := financing.NewID()
		invoiceID1 := financing.NewID()
		invoiceID2 := financing.NewID()

		s.publisEvents(
			financing.NewIssuerCreatedEvent(issuerID),
			financing.NewInvoiceCreatedEvent(invoiceID1, issuerID, 15),
			financing.NewInvoiceCreatedEvent(invoiceID2, issuerID, 20),
		)

		expected := []*pb.Invoice{
			{
				Id:          grpc.ConvertID(invoiceID1),
				IssuerId:    grpc.ConvertID(issuerID),
				AskingPrice: &pb.Money{Amount: 15},
				Status:      pb.InvoiceStatus_AVAILABLE,
				WinningBid:  nil,
			},
			{
				Id:          grpc.ConvertID(invoiceID2),
				IssuerId:    grpc.ConvertID(issuerID),
				AskingPrice: &pb.Money{Amount: 30},
				Status:      pb.InvoiceStatus_AVAILABLE,
				WinningBid:  nil,
			},
		}

		invoices, err := s.queries.AllInvoices(context.Background(), &emptypb.Empty{})
		require.NoError(t, err)

		require.Equal(t, expected, invoices)
	})

	t.Run(`
	GIVEN there is an issuer registered 
	AND she sells an invoice with asking price of 15
	AND there is an investor registered with 20 balance
	AND that investor bids on the same invoice for 20
	WHEN all invoices are queried
	THEN the invoice is obtained with status financed 
	AND the winning bid of that investor with an amount of 20`, func(t *testing.T) {
		issuerID := financing.NewID()
		invoiceID := financing.NewID()
		investorID := financing.NewID()

		s.publisEvents(
			financing.NewIssuerCreatedEvent(issuerID),
			financing.NewInvoiceCreatedEvent(invoiceID, issuerID, 15),
			financing.NewInvestorCreatedEvent(investorID),
			financing.NewInvestorFundsAddedEvent(investorID, 20),
			financing.NewBidOnInvoicePlacedEvent(investorID, investorID, 20),
		)

		expected := []*pb.Invoice{
			{
				Id:          grpc.ConvertID(invoiceID),
				IssuerId:    grpc.ConvertID(issuerID),
				AskingPrice: &pb.Money{Amount: 15},
				Status:      pb.InvoiceStatus_FINANCED,
				WinningBid: &pb.Bid{
					InvestorId: grpc.ConvertID(investorID),
					Amount:     &pb.Money{Amount: 20},
				},
			},
		}

		invoices, err := s.queries.AllInvoices(context.Background(), &emptypb.Empty{})
		require.NoError(t, err)

		require.Equal(t, expected, invoices)
	})

	t.Run(`
	GIVEN there is an issuer registered 
	AND she sells an invoice with asking price of 15
	AND there is an investor registered with 20 balance
	AND that investor bids on the same invoice for 20
	AND the issuer accepts the financing
	WHEN all invoices are queried
	THEN the invoice is obtained with status accepted 
	AND the winning bid of that investor with an amount of 20`, func(t *testing.T) {
		issuerID := financing.NewID()
		invoiceID := financing.NewID()
		investorID := financing.NewID()

		s.publisEvents(
			financing.NewIssuerCreatedEvent(issuerID),
			financing.NewInvoiceCreatedEvent(invoiceID, issuerID, 15),
			financing.NewInvestorCreatedEvent(investorID),
			financing.NewInvestorFundsAddedEvent(investorID, 20),
			financing.NewBidOnInvoicePlacedEvent(investorID, investorID, 20),
			financing.NewInvoiceApprovedEvent(invoiceID, 15, financing.NewBid(investorID, 20)),
		)

		expected := []*pb.Invoice{
			{
				Id:          grpc.ConvertID(invoiceID),
				IssuerId:    grpc.ConvertID(issuerID),
				AskingPrice: &pb.Money{Amount: 15},
				Status:      pb.InvoiceStatus_APPROVED,
				WinningBid: &pb.Bid{
					InvestorId: grpc.ConvertID(investorID),
					Amount:     &pb.Money{Amount: 20},
				},
			},
		}

		invoices, err := s.queries.AllInvoices(context.Background(), &emptypb.Empty{})
		require.NoError(t, err)

		require.Equal(t, expected, invoices)
	})

	t.Run(`
	GIVEN there is an issuer registered 
	AND she sells an invoice with asking price of 15
	AND there is an investor registered with 20 balance
	AND that investor bids on the same invoice for 20
	AND the issuer reverses the financing
	WHEN all invoices are queried
	THEN the invoice is obtained with status reversed 
	AND the winning bid of that investor with an amount of 20`, func(t *testing.T) {
		issuerID := financing.NewID()
		invoiceID := financing.NewID()
		investorID := financing.NewID()

		s.publisEvents(
			financing.NewIssuerCreatedEvent(issuerID),
			financing.NewInvoiceCreatedEvent(invoiceID, issuerID, 15),
			financing.NewInvestorCreatedEvent(investorID),
			financing.NewInvestorFundsAddedEvent(investorID, 20),
			financing.NewBidOnInvoicePlacedEvent(investorID, investorID, 20),
			financing.NewInvoiceReversedEvent(invoiceID, 15, financing.NewBid(investorID, 20)),
		)

		expected := []*pb.Invoice{
			{
				Id:          grpc.ConvertID(invoiceID),
				IssuerId:    grpc.ConvertID(issuerID),
				AskingPrice: &pb.Money{Amount: 15},
				Status:      pb.InvoiceStatus_REVERSED,
				WinningBid: &pb.Bid{
					InvestorId: grpc.ConvertID(investorID),
					Amount:     &pb.Money{Amount: 20},
				},
			},
		}

		invoices, err := s.queries.AllInvoices(context.Background(), &emptypb.Empty{})
		require.NoError(t, err)

		require.Equal(t, expected, invoices)
	})
}
