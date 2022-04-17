package grpc

import (
	"context"
	"fmt"
	"net"

	"github.com/pperaltaisern/financing/pkg/grpc/pb"
	"github.com/pperaltaisern/financing/pkg/query"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
)

type QueryServer struct {
	network, address string
	server           *grpc.Server
	investors        query.InvestorQueries
	invoices         query.InvoiceQueries
	issuers          query.IssuerQueries
	pb.UnimplementedQueriesServer
}

func NewQueryServer(network, address string, investors query.InvestorQueries, invoices query.InvoiceQueries, issuers query.IssuerQueries) *QueryServer {
	s := &QueryServer{
		network:   network,
		address:   address,
		server:    grpc.NewServer(),
		investors: investors,
		invoices:  invoices,
		issuers:   issuers,
	}
	pb.RegisterQueriesServer(s.server, s)
	reflection.Register(s.server)
	return s
}

func (s *QueryServer) Open() error {
	l, err := net.Listen(s.network, s.address)
	if err != nil {
		return fmt.Errorf("opening QueryServer for network: %s and address: %s, err: %v", s.network, s.address, err)
	}
	return s.server.Serve(l)
}

func (s *QueryServer) Close() {
	s.server.GracefulStop()
}

func (s *QueryServer) AllInvestors(context.Context, *emptypb.Empty) (*pb.AllInvestorsReply, error) {
	investors, err := s.investors.All()
	if err != nil {
		return nil, err
	}
	pbInvestors := make([]*pb.Investor, len(investors))
	for i, investor := range investors {
		pbInvestors[i] = &pb.Investor{
			Id:       ConvertID(investor.ID),
			Balance:  ConvertMoney(investor.Balance),
			Reserved: ConvertMoney(investor.Reserved),
		}
	}

	return &pb.AllInvestorsReply{
		Investors: pbInvestors,
	}, nil
}

func (s *QueryServer) AllInvoices(context.Context, *emptypb.Empty) (*pb.AllInvoicesReply, error) {
	invoices, err := s.invoices.All()
	if err != nil {
		return nil, err
	}
	pbInvoices := make([]*pb.Invoice, len(invoices))
	for i, invoice := range invoices {
		var winningBid *pb.Bid
		if invoice.WinningBid != nil {
			winningBid = &pb.Bid{
				InvestorId: ConvertID(invoice.WinningBid.InvestorID),
				Amount:     ConvertMoney(invoice.WinningBid.Amount),
			}
		}

		pbInvoices[i] = &pb.Invoice{
			Id:          ConvertID(invoice.ID),
			IssuerId:    ConvertID(invoice.IssuerID),
			AskingPrice: ConvertMoney(invoice.AskingPrice),
			Status:      convertInvoiceStatus(invoice.Status),
			WinningBid:  winningBid,
		}
	}

	return &pb.AllInvoicesReply{
		Invoices: pbInvoices,
	}, nil
}

func (s *QueryServer) AllIssuers(context.Context, *emptypb.Empty) (*pb.AllIssuersReply, error) {
	issuers, err := s.issuers.All()
	if err != nil {
		return nil, err
	}
	pbIssuers := make([]*pb.Issuer, len(issuers))
	for i, issuer := range issuers {
		pbIssuers[i] = &pb.Issuer{
			Id: ConvertID(issuer.ID),
		}
	}

	return &pb.AllIssuersReply{
		Issuers: pbIssuers,
	}, nil
}
