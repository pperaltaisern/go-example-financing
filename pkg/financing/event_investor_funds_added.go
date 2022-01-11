package financing

type InvestorFundsAddedEvent struct {
	InvestorID ID `json:"-"`
	Amount     Money
}

func NewInvestorFundsAddedEvent(investorID ID, amount Money) *InvestorFundsAddedEvent {
	return &InvestorFundsAddedEvent{
		InvestorID: investorID,
		Amount:     amount,
	}
}

func (e *InvestorFundsAddedEvent) EventName() string { return "InvestorFundsAddedEvent" }

func (e *InvestorFundsAddedEvent) WithAggregateID(id string) {
	e.InvestorID = NewIDFromString(id)
}
