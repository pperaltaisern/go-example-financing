package projection

import "github.com/pperaltaisern/financing/pkg/financing"

type EventProjector interface {
	ProjectInvestorCreatedEvent(*financing.InvestorCreatedEvent) error
	ProjectInvestorFundsAddedEvent(*financing.InvestorFundsAddedEvent) error
	ProjectBidOnInvoicePlacedEvent(*financing.BidOnInvoicePlacedEvent) error
	ProjectInvestorFundsReleasedEvent(*financing.InvestorFundsReleasedEvent) error
	ProjectInvestorFundsCommittedEvent(*financing.InvestorFundsCommittedEvent) error
}

var _ EventProjector = (*MockEventProjector)(nil)

type MockEventProjector struct {
	ProjectInvestorCreatedEventFn        func(*financing.InvestorCreatedEvent) error
	ProjectInvestorFundsAddedEventFn     func(*financing.InvestorFundsAddedEvent) error
	ProjectBidOnInvoicePlacedEventFn     func(*financing.BidOnInvoicePlacedEvent) error
	ProjectInvestorFundsReleasedEventFn  func(*financing.InvestorFundsReleasedEvent) error
	ProjectInvestorFundsCommittedEventFn func(*financing.InvestorFundsCommittedEvent) error
}

func (m *MockEventProjector) ProjectInvestorCreatedEvent(e *financing.InvestorCreatedEvent) error {
	return m.ProjectInvestorCreatedEventFn(e)
}
func (m *MockEventProjector) ProjectInvestorFundsAddedEvent(e *financing.InvestorFundsAddedEvent) error {
	return m.ProjectInvestorFundsAddedEventFn(e)
}
func (m *MockEventProjector) ProjectBidOnInvoicePlacedEvent(e *financing.BidOnInvoicePlacedEvent) error {
	return m.ProjectBidOnInvoicePlacedEventFn(e)
}
func (m *MockEventProjector) ProjectInvestorFundsReleasedEvent(e *financing.InvestorFundsReleasedEvent) error {
	return m.ProjectInvestorFundsReleasedEventFn(e)
}
func (m *MockEventProjector) ProjectInvestorFundsCommittedEvent(e *financing.InvestorFundsCommittedEvent) error {
	return m.ProjectInvestorFundsCommittedEventFn(e)
}
