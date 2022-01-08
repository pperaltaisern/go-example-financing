package financing

import (
	"context"
	"ledger/internal/esrc"
)

type IssuerRepository interface {
	Contains(context.Context, ID) (bool, error)
	Add(context.Context, *Issuer) error
}

func NewIssuerRepository(es esrc.EventStore) IssuerRepository {
	return issuerRepository{es: es}
}

type issuerRepository struct {
	es esrc.EventStore
}

func (r issuerRepository) Contains(ctx context.Context, id ID) (bool, error) {
	return r.es.Contains(ctx, id)
}

func (i issuerRepository) Add(ctx context.Context, iss *Issuer) error {
	rawEvents, err := esrc.MarshalEventsJSON(iss.aggregate.Events())
	if err != nil {
		return err
	}

	const aggregateType esrc.AggregateType = "issuer"
	return i.es.Create(ctx, aggregateType, iss.id, rawEvents)
}
