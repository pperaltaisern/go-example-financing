package command

import "ledger/pkg/financing"

type ApproveFinancing struct {
	InvoiceID financing.ID
}
