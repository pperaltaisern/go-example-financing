package config

import (
	"github.com/pperaltaisern/financing/pkg/financing"
)

type Repositories struct {
	Issuers   financing.IssuerRepository
	Investors financing.InvestorRepository
	Invoices  financing.InvoiceRepository
}
