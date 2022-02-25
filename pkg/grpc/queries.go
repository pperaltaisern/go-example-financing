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
	pb.UnimplementedQueriesServer
}

func NewQueryServer(network, address string, investors query.InvestorQueries, invoices query.InvoiceQueries) *QueryServer {
	s := &QueryServer{
		network:   network,
		address:   address,
		server:    grpc.NewServer(),
		investors: investors,
		invoices:  invoices,
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
	return nil, nil
}

func (s *QueryServer) AllInvoices(context.Context, *emptypb.Empty) (*pb.AllInvoicesReply, error) {
	return nil, nil
}
