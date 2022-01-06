package financing

type InvestorFundsReleasedEvent struct {
	InvestorID ID
	Amount     Money
}

func NewInvestorFundsReleasedEvent(investorID ID, amount Money) InvestorFundsReleasedEvent {
	return InvestorFundsReleasedEvent{
		InvestorID: investorID,
		Amount:     amount,
	}
}

func (e InvestorFundsReleasedEvent) IsEvent() {}

func (e InvestorFundsReleasedEvent) EventName() string { return "InvestorFundsReleasedEvent" }
