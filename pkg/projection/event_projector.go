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
