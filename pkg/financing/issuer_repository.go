package financing

import (
	"context"

	"github.com/pperaltaisern/financing/internal/esrc"
)

type IssuerRepository interface {
	Contains(context.Context, ID) (bool, error)
	Add(context.Context, *Issuer) error
}

type issuerRepository struct {
	r *esrc.Repository[*Issuer]
}

func NewIssuerRepository(es esrc.EventStore) IssuerRepository {
	return issuerRepository{
		r: esrc.NewRepository[*Issuer](es, nil, nil),
	}
}

func (r issuerRepository) Contains(ctx context.Context, id ID) (bool, error) {
	return r.r.Contains(ctx, id)
}

func (r issuerRepository) Add(ctx context.Context, iss *Issuer) error {
	return r.r.Add(ctx, iss)
}
