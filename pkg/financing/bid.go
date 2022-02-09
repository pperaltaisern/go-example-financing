package financing

type Bid struct {
	InvestorID ID
	Amount     Money
}

func NewBid(investorID ID, amount Money) Bid {
	return Bid{
		InvestorID: investorID,
		Amount:     amount,
	}
}
