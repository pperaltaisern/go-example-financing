package financing

type InvestorFundsReleasedEvent struct {
	InvestorID ID `json:"-"`
	Amount     Money
}

func NewInvestorFundsReleasedEvent(investorID ID, amount Money) *InvestorFundsReleasedEvent {
	return &InvestorFundsReleasedEvent{
		InvestorID: investorID,
		Amount:     amount,
	}
}

func (e *InvestorFundsReleasedEvent) EventName() string { return "InvestorFundsReleasedEvent" }

func (e *InvestorFundsReleasedEvent) WithAggregateID(id string) {
	e.InvestorID = NewIDFrom(id)
}
