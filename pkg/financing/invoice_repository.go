package financing

import "context"

type InvoiceRepository interface {
	ByID(context.Context, ID) (*Invoice, error)
	Update(context.Context, *Invoice) error
	Add(context.Context, *Invoice) error
}
