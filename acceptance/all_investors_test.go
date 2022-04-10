package acceptance

import (
	"context"
	"testing"

	"github.com/pperaltaisern/financing/acceptance/builder"
	"github.com/pperaltaisern/financing/pkg/financing"
	"github.com/pperaltaisern/financing/pkg/grpc"
	"github.com/pperaltaisern/financing/pkg/grpc/pb"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *QueriesSuite) TestAllInvestors() {
	t := s.T()

	t.Run(`
	GIVEN that there isn't any investor registered
	WHEN all investors are queried
	THEN no results are obtained`, func(t *testing.T) {
		reply, err := s.queries.AllInvestors(context.Background(), &emptypb.Empty{})
		require.NoError(t, err)

		require.Empty(t, reply.Investors)
	})

	t.Run(`
	GIVEN that we register 2 investors
	WHEN all investors are queried
	THEN we obtain those 2 investors with empty balance`, func(t *testing.T) {
		id1 := financing.NewID()
		id2 := financing.NewID()

		s.publisEvents(
			s.newRelayEvent(id1, financing.NewInvestorCreatedEvent(id1)),
			s.newRelayEvent(id2, financing.NewInvestorCreatedEvent(id2)),
		)

		expected := []*pb.Investor{
			builder.NewInvestor(id1).Build(),
			builder.NewInvestor(id2).Build(),
		}

		reply, err := s.queries.AllInvestors(context.Background(), &emptypb.Empty{})
		require.NoError(t, err)

		RequireJSONEq(t, expected, reply.Investors)
	})

	t.Run(`
	GIVEN that we have a registered investor
	AND 20 funds are added
	AND 30 funds are added
	WHEN all investors are queried
	THEN we find that investor with 50 balance`, func(t *testing.T) {
		id := financing.NewID()

		s.publisEvents(
			s.newRelayEvent(id, financing.NewInvestorCreatedEvent(id)),
			s.newRelayEvent(id, financing.NewInvestorFundsAddedEvent(id, 20)),
			s.newRelayEvent(id, financing.NewInvestorFundsAddedEvent(id, 30)),
		)

		expected := builder.NewInvestor(id).WithBalance(50).Build()

		reply, err := s.queries.AllInvestors(context.Background(), &emptypb.Empty{})
		require.NoError(t, err)

		RequireJSONEq(t, expected, findInvestor(reply.Investors, id))
	})

	t.Run(`
	GIVEN that we have a registered investor
	AND 30 funds are added
	AND an Issuer registered
	AND that issuer sells an invoice with an asking price of 15
	AND the investors bids on the invoice for 20
	WHEN all investors are queried
	THEN we find that investor with 15 balance and 15 reserved funds`, func(t *testing.T) {
		investorID := financing.NewID()
		issuerID := financing.NewID()
		invoiceID := financing.NewID()

		s.publisEvents(
			s.newRelayEvent(investorID, financing.NewInvestorCreatedEvent(investorID)),
			s.newRelayEvent(investorID, financing.NewInvestorFundsAddedEvent(investorID, 30)),
			s.newRelayEvent(issuerID, financing.NewIssuerCreatedEvent(issuerID)),
			s.newRelayEvent(invoiceID, financing.NewInvoiceCreatedEvent(invoiceID, issuerID, 15)),
			s.newRelayEvent(investorID, financing.NewBidOnInvoicePlacedEvent(investorID, invoiceID, 20)),
			s.newRelayEvent(investorID, financing.NewInvestorFundsReleasedEvent(investorID, 5)),
		)

		expected := builder.NewInvestor(investorID).WithBalance(15).WithReserved(15).Build()

		reply, err := s.queries.AllInvestors(context.Background(), &emptypb.Empty{})
		require.NoError(t, err)

		RequireJSONEq(t, expected, findInvestor(reply.Investors, investorID))
	})

	t.Run(`
	GIVEN an invoice with asking price of 15 is financed by an investor
	AND the investor has 0 balance and 15 reserved
	AND the issuer approves the financing
	WHEN all investors are queried
	THEN we find that investor with 0 balance and 0 reserved funds`, func(t *testing.T) {
		investorID := financing.NewID()
		issuerID := financing.NewID()
		invoiceID := financing.NewID()

		s.publisEvents(
			s.newRelayEvent(investorID, financing.NewInvestorCreatedEvent(investorID)),
			s.newRelayEvent(investorID, financing.NewInvestorFundsAddedEvent(investorID, 15)),
			s.newRelayEvent(issuerID, financing.NewIssuerCreatedEvent(issuerID)),
			s.newRelayEvent(invoiceID, financing.NewInvoiceCreatedEvent(invoiceID, issuerID, 15)),
			s.newRelayEvent(investorID, financing.NewBidOnInvoicePlacedEvent(investorID, invoiceID, 15)),
			s.newRelayEvent(invoiceID, financing.NewInvoiceApprovedEvent(invoiceID, 15, financing.NewBid(investorID, 15))),
			s.newRelayEvent(investorID, financing.NewInvestorFundsCommittedEvent(investorID, 15)),
		)

		expected := builder.NewInvestor(investorID).Build()

		reply, err := s.queries.AllInvestors(context.Background(), &emptypb.Empty{})
		require.NoError(t, err)

		RequireJSONEq(t, expected, findInvestor(reply.Investors, investorID))
	})

	t.Run(`
	GIVEN an invoice with asking price of 15 is financed by an investor
	AND the investor has 0 balance and 15 reserved
	AND the issuer reverses the financing
	WHEN all investors are queried
	THEN we find that investor with 15 balance and 0 reserved funds`, func(t *testing.T) {
		investorID := financing.NewID()
		issuerID := financing.NewID()
		invoiceID := financing.NewID()

		s.publisEvents(
			s.newRelayEvent(investorID, financing.NewInvestorCreatedEvent(investorID)),
			s.newRelayEvent(investorID, financing.NewInvestorFundsAddedEvent(investorID, 15)),
			s.newRelayEvent(issuerID, financing.NewIssuerCreatedEvent(issuerID)),
			s.newRelayEvent(invoiceID, financing.NewInvoiceCreatedEvent(invoiceID, issuerID, 15)),
			s.newRelayEvent(investorID, financing.NewBidOnInvoicePlacedEvent(investorID, invoiceID, 15)),
			s.newRelayEvent(invoiceID, financing.NewInvoiceReversedEvent(invoiceID, 15, financing.NewBid(investorID, 15))),
			s.newRelayEvent(investorID, financing.NewInvestorFundsReleasedEvent(invoiceID, 15)),
		)

		expected := builder.NewInvestor(investorID).WithBalance(15).Build()

		reply, err := s.queries.AllInvestors(context.Background(), &emptypb.Empty{})
		require.NoError(t, err)

		RequireJSONEq(t, expected, findInvestor(reply.Investors, investorID))
	})
}

func findInvestor(investors []*pb.Investor, id financing.ID) *pb.Investor {
	for _, investor := range investors {
		if investor.Id.Value == grpc.ConvertID(id).Value {
			return investor
		}
	}
	return nil
}
