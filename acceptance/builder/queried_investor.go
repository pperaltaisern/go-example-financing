package builder

import (
	"github.com/pperaltaisern/financing/pkg/financing"
	"github.com/pperaltaisern/financing/pkg/grpc"
	"github.com/pperaltaisern/financing/pkg/grpc/pb"
)

type Investor struct {
	id                financing.ID
	balance, reserved float64
}

func NewInvestor(id financing.ID) *Investor {
	return &Investor{
		id: id,
	}
}

func (inv *Investor) WithBalance(m float64) *Investor {
	inv.balance = m
	return inv
}

func (inv *Investor) WithReserved(m float64) *Investor {
	inv.reserved = m
	return inv
}

func (inv *Investor) Build() *pb.Investor {
	return &pb.Investor{
		Id:       grpc.ConvertID(inv.id),
		Balance:  grpc.ConvertMoney(financing.Money(inv.balance)),
		Reserved: grpc.ConvertMoney(financing.Money(inv.reserved)),
	}
}
