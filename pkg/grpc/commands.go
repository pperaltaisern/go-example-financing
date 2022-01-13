package grpc

import (
	"context"
	"fmt"
	"net"

	"github.com/pperaltaisern/financing/pkg/command"
	"github.com/pperaltaisern/financing/pkg/financing"
	"github.com/pperaltaisern/financing/pkg/grpc/pb"

	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
)

type CommandServer struct {
	network, address string
	commandBus       *cqrs.CommandBus
	server           *grpc.Server

	pb.UnimplementedCommandsServer
}

func NewCommandServer(network, address string, bus *cqrs.CommandBus) *CommandServer {
	s := &CommandServer{
		commandBus: bus,
		server:     grpc.NewServer(),
	}
	pb.RegisterCommandsServer(s.server, s)
	reflection.Register(s.server)
	return s
}

func (s *CommandServer) Open() error {
	l, err := net.Listen(s.network, s.address)
	if err != nil {
		return err
	}
	return s.server.Serve(l)
}

func (s *CommandServer) Close() {
	s.server.GracefulStop()
}

func (s *CommandServer) SellInvoice(ctx context.Context, pbcmd *pb.SellInvoiceCommand) (*pb.UUID, error) {
	if pbcmd.AskingPrice.Amount <= 0 {
		return nil, fmt.Errorf("asking price must be higher than 0")
	}
	issuerID, err := financing.TryNewIDFromString(pbcmd.IssuerId.Value)
	if err != nil {
		return nil, fmt.Errorf("issuer id %s is not valid", pbcmd.IssuerId.Value)
	}

	invoiceID := financing.NewID()
	cmd := command.NewSellInvoice(
		invoiceID,
		issuerID,
		financing.Money(pbcmd.AskingPrice.Amount))

	err = s.commandBus.Send(ctx, cmd)
	if err != nil {
		return nil, err
	}
	return convertID(invoiceID), nil
}

func (s *CommandServer) BidOnInvoice(ctx context.Context, pbcmd *pb.BidOnInvoiceCommand) (*emptypb.Empty, error) {
	if pbcmd.Bid.Amount <= 0 {
		return nil, fmt.Errorf("bid amount must be higher than 0")
	}

	investorID, err := financing.TryNewIDFromString(pbcmd.InvestorId.Value)
	if err != nil {
		return nil, fmt.Errorf("investor id %s is not valid", pbcmd.InvestorId.Value)
	}
	invoiceID, err := financing.TryNewIDFromString(pbcmd.InvoiceId.Value)
	if err != nil {
		return nil, fmt.Errorf("invoice id %s is not valid", pbcmd.InvoiceId.Value)
	}

	cmd := command.NewBidOnInvoice(
		investorID,
		invoiceID,
		financing.Money(pbcmd.Bid.Amount))

	err = s.commandBus.Send(ctx, cmd)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *CommandServer) ApproveFinancing(ctx context.Context, pbcmd *pb.ApproveFinancingCommand) (*emptypb.Empty, error) {
	invoiceID, err := financing.TryNewIDFromString(pbcmd.InvoiceId.Value)
	if err != nil {
		return nil, fmt.Errorf("invoice id %s is not valid", pbcmd.InvoiceId.Value)
	}

	cmd := command.NewApproveFinancing(invoiceID)
	err = s.commandBus.Send(ctx, cmd)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
func (s *CommandServer) ReverseFinancing(ctx context.Context, pbcmd *pb.ReverseFinancingCommand) (*emptypb.Empty, error) {
	invoiceID, err := financing.TryNewIDFromString(pbcmd.InvoiceId.Value)
	if err != nil {
		return nil, fmt.Errorf("invoice id %s is not valid", pbcmd.InvoiceId.Value)
	}

	cmd := command.NewReverseFinancing(invoiceID)
	err = s.commandBus.Send(ctx, cmd)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
