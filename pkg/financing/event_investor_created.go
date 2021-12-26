package financing

type InvestorCreatedEvent struct {
	InvestorID ID
}

func NewInvestorCreatedEvent(investorID ID) InvestorCreatedEvent {
	return InvestorCreatedEvent{
		InvestorID: investorID,
	}
}

func (e InvestorCreatedEvent) IsEvent() {}

func (e InvestorCreatedEvent) Name() string { return "InvestorCreatedEvent" }
