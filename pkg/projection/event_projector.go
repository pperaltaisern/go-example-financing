package projection

import "github.com/pperaltaisern/financing/pkg/financing"

type EventProjector interface {
	ProjectInvestorCreatedEvent(*financing.InvestorCreatedEvent) error
	ProjectInvestorFundsAddedEvent(*financing.InvestorFundsAddedEvent) error
	ProjectBidOnInvoicePlacedEvent(*financing.BidOnInvoicePlacedEvent) error
	ProjectInvestorFundsReleasedEvent(*financing.InvestorFundsReleasedEvent) error
	ProjectInvestorFundsCommittedEvent(*financing.InvestorFundsCommittedEvent) error
	ProjectInvoiceCreatedEvent(*financing.InvoiceCreatedEvent) error
	ProjectInvoiceFinancedEvent(*financing.InvoiceFinancedEvent) error
	ProjectInvoiceReversedEvent(*financing.InvoiceReversedEvent) error
	ProjectInvoiceApprovedEvent(*financing.InvoiceApprovedEvent) error
	ProjectIssuerCreatedEvent(*financing.IssuerCreatedEvent) error
}

var _ EventProjector = (*MockEventProjector)(nil)

type MockEventProjector struct {
	ProjectInvestorCreatedEventFn        func(*financing.InvestorCreatedEvent) error
	ProjectInvestorFundsAddedEventFn     func(*financing.InvestorFundsAddedEvent) error
	ProjectBidOnInvoicePlacedEventFn     func(*financing.BidOnInvoicePlacedEvent) error
	ProjectInvestorFundsReleasedEventFn  func(*financing.InvestorFundsReleasedEvent) error
	ProjectInvestorFundsCommittedEventFn func(*financing.InvestorFundsCommittedEvent) error
	ProjectInvoiceCreatedEventFn         func(*financing.InvoiceCreatedEvent) error
	ProjectInvoiceFinancedEventFn        func(*financing.InvoiceFinancedEvent) error
	ProjectInvoiceReversedEventFn        func(*financing.InvoiceReversedEvent) error
	ProjectInvoiceApprovedEventFn        func(*financing.InvoiceApprovedEvent) error
	ProjectIssuerCreatedEventFn          func(*financing.IssuerCreatedEvent) error
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
func (m *MockEventProjector) ProjectInvoiceCreatedEvent(e *financing.InvoiceCreatedEvent) error {
	return m.ProjectInvoiceCreatedEventFn(e)
}
func (m *MockEventProjector) ProjectInvoiceFinancedEvent(e *financing.InvoiceFinancedEvent) error {
	return m.ProjectInvoiceFinancedEventFn(e)
}
func (m *MockEventProjector) ProjectInvoiceReversedEvent(e *financing.InvoiceReversedEvent) error {
	return m.ProjectInvoiceReversedEventFn(e)
}
func (m *MockEventProjector) ProjectInvoiceApprovedEvent(e *financing.InvoiceApprovedEvent) error {
	return m.ProjectInvoiceApprovedEventFn(e)
}
func (m *MockEventProjector) ProjectIssuerCreatedEvent(e *financing.IssuerCreatedEvent) error {
	return m.ProjectIssuerCreatedEventFn(e)
}
