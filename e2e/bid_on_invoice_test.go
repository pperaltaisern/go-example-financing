package e2e

import (
	"context"
	"testing"

	"github.com/pperaltaisern/financing/pkg/financing"
	"github.com/pperaltaisern/financing/pkg/grpc"
	"github.com/pperaltaisern/financing/pkg/grpc/pb"
	"github.com/stretchr/testify/require"
)

func (s *Suite) TestBidOnInvoice() {
	t := s.T()

	t.Run(`
	GIVEN an investor has been registered 
	AND he has up to 100 balance 
	AND there is an invoice with asking price of 10
	WHEN he bids on an invoice for 20
	THEN 20 funds are reserved from the investor account
	THEN invoice is set to financed
	THEN 10 funds are released from the investor account`, func(t *testing.T) {
		s.FinanceAnInvoice()
	})

	t.Run(`
	GIVEN an investor has been registered 
	AND he has up to 100 balance 
	AND there is an invoice with asking price of 10
	WHEN he bids on an invoice for 5
	THEN 5 funds are reserved from the investor account
	THEN bid on invoice is rejected
	THEN 5 funds are released from the investor account`, func(t *testing.T) {

		askingPrice := financing.Money(10)
		invoiceID := s.RegisterIssuerAndSellInvoice(askingPrice)
		investorID := financing.NewID()
		s.RegisterInvestor(investorID, 100)

		bidAmount := financing.Money(5)
		cmd := &pb.BidOnInvoiceCommand{
			InvestorId: grpc.ConvertID(investorID),
			InvoiceId:  grpc.ConvertID(invoiceID),
			Bid: &pb.Money{
				Amount: float64(bidAmount),
			},
		}
		_, err := s.commands.BidOnInvoice(context.Background(), cmd)
		require.NoError(t, err)

		eventAssertions := []EventAssertion{
			{
				Expected: financing.NewBidOnInvoicePlacedEvent(investorID, invoiceID, bidAmount),
				Actual:   &financing.BidOnInvoicePlacedEvent{},
			},
			{
				Expected: financing.NewBidOnInvoiceRejectedEvent(
					invoiceID,
					financing.Bid{
						InvestorID: investorID,
						Amount:     bidAmount,
					}),
				Actual: &financing.BidOnInvoiceRejectedEvent{},
			},
			{
				Expected: financing.NewInvestorFundsReleasedEvent(investorID, bidAmount),
				Actual:   &financing.InvestorFundsReleasedEvent{},
			},
		}

		s.expectEvents(t, eventAssertions...)
	})

	t.Run(`
	GIVEN an investor has been registered 
	AND he has up to 5 balance 
	AND there is an invoice with asking price of 10
	WHEN he bids on an invoice for 10
	THEN no funds are reserved since he hasn't enough balance`, func(t *testing.T) {

		askingPrice := financing.Money(10)
		invoiceID := s.RegisterIssuerAndSellInvoice(askingPrice)
		investorID := financing.NewID()
		s.RegisterInvestor(investorID, 5)

		bidAmount := financing.Money(10)
		cmd := &pb.BidOnInvoiceCommand{
			InvestorId: grpc.ConvertID(investorID),
			InvoiceId:  grpc.ConvertID(invoiceID),
			Bid: &pb.Money{
				Amount: float64(bidAmount),
			},
		}
		_, err := s.commands.BidOnInvoice(context.Background(), cmd)
		require.NoError(t, err)
	})
}

type FinanceAnInvoiceResult struct {
	InvoiceID   financing.ID
	InvoiceCost financing.Money
	Bid         financing.Bid
}

func (s *Suite) FinanceAnInvoice() FinanceAnInvoiceResult {
	t := s.T()

	askingPrice := financing.Money(10)
	invoiceID := s.RegisterIssuerAndSellInvoice(askingPrice)
	investorID := financing.NewID()
	s.RegisterInvestor(investorID, 100)

	bidAmount := financing.Money(20)
	cmd := &pb.BidOnInvoiceCommand{
		InvestorId: grpc.ConvertID(investorID),
		InvoiceId:  grpc.ConvertID(invoiceID),
		Bid: &pb.Money{
			Amount: float64(bidAmount),
		},
	}
	_, err := s.commands.BidOnInvoice(context.Background(), cmd)
	require.NoError(t, err)

	eventAssertions := []EventAssertion{
		{
			Expected: financing.NewBidOnInvoicePlacedEvent(investorID, invoiceID, bidAmount),
			Actual:   &financing.BidOnInvoicePlacedEvent{},
		},
		{
			Expected: financing.NewInvoiceFinancedEvent(
				invoiceID,
				askingPrice,
				financing.Bid{
					InvestorID: investorID,
					Amount:     bidAmount,
				}),
			Actual: &financing.InvoiceFinancedEvent{},
		},
		{
			Expected: financing.NewInvestorFundsReleasedEvent(investorID, bidAmount-askingPrice),
			Actual:   &financing.InvestorFundsReleasedEvent{},
		},
	}

	s.expectEvents(t, eventAssertions...)

	return FinanceAnInvoiceResult{
		InvoiceID:   invoiceID,
		InvoiceCost: askingPrice,
		Bid: financing.Bid{
			InvestorID: investorID,
			Amount:     bidAmount,
		},
	}
}
