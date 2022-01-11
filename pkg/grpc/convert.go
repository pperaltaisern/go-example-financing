package grpc

import (
	"ledger/pkg/financing"
	"ledger/pkg/grpc/pb"
)

func convertID(id financing.ID) *pb.UUID {
	return &pb.UUID{
		Value: id.String(),
	}
}
