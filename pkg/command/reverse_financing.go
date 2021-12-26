package command

import "ledger/pkg/financing"

type ReverseFinancing struct {
	InvoiceID financing.ID
}
