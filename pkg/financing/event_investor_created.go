package financing

type InvestorCreatedEvent struct {
	InvestorID ID
}

func NewInvestorCreatedEvent(investorID ID) *InvestorCreatedEvent {
	return &InvestorCreatedEvent{
		InvestorID: investorID,
	}
}

func (e *InvestorCreatedEvent) EventName() string { return "InvestorCreatedEvent" }

func (e *InvestorCreatedEvent) WithAggregateID(id string) {
	e.InvestorID = NewIDFrom(id)
}
