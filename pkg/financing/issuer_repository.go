package financing

import "context"

type IssuerRepository interface {
	Contains(context.Context, ID) (bool, error)
	Add(context.Context, *Issuer) error
}
