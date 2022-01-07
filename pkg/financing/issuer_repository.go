package financing

import (
	"context"
	"encoding/json"
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

func (i issuerRepository) Add(ctx context.Context, inv *Issuer) error {
	rawEvents := make([]esrc.RawEvent, len(inv.aggregate.Events()))
	for i, e := range inv.aggregate.Events() {
		b, err := json.Marshal(e)
		if err != nil {
			return err
		}
		rawEvents[i] = esrc.RawEvent{Name: e.EventName(), Data: b}
	}

	const aggregateType esrc.AggregateType = "issuer"
	return i.es.Create(ctx, aggregateType, inv.id, rawEvents)
}
