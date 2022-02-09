package grpc

import (
	"github.com/pperaltaisern/financing/pkg/financing"
	"github.com/pperaltaisern/financing/pkg/grpc/pb"
)

func ConvertID(id financing.ID) *pb.UUID {
	return &pb.UUID{
		Value: id.String(),
	}
}
