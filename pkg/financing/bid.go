package financing

type Bid struct {
	IvestorID ID
	Amount    Money
}

func NewBid(investorID ID, amount Money) Bid {
	return Bid{
		IvestorID: investorID,
		Amount:    amount,
	}
}
