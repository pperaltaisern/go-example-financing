package grpc

import (
	"github.com/pperaltaisern/financing/pkg/financing"
	"github.com/pperaltaisern/financing/pkg/grpc/pb"
	"github.com/pperaltaisern/financing/pkg/query"
)

func ConvertID(id financing.ID) *pb.UUID {
	return &pb.UUID{
		Value: id.String(),
	}
}

func ConvertMoney(money financing.Money) *pb.Money {
	return &pb.Money{
		Amount: float64(money),
	}
}

func convertInvoiceStatus(status query.InvoiceStatus) pb.InvoiceStatus {
	switch status {
	case query.InvoiceStatusAvailable:
		return pb.InvoiceStatus_AVAILABLE
	case query.InvoiceStatusFinanced:
		return pb.InvoiceStatus_FINANCED
	case query.InvoiceStatusApproved:
		return pb.InvoiceStatus_APPROVED
	case query.InvoiceStatusReversed:
		return pb.InvoiceStatus_REVERSED
	}
	panic("unknown invoice status")
}
