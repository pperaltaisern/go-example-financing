package financing

type InvestorFundsCommittedEvent struct {
	InvestorID ID `json:"-"`
	Amount     Money
}

func NewInvestorFundsCommittedEvent(investorID ID, amount Money) *InvestorFundsCommittedEvent {
	return &InvestorFundsCommittedEvent{
		InvestorID: investorID,
		Amount:     amount,
	}
}

func (e *InvestorFundsCommittedEvent) EventName() string { return "InvestorFundsCommittedEvent" }

func (e *InvestorFundsCommittedEvent) WithAggregateID(id string) {
	e.InvestorID = NewIDFromString(id)
}
