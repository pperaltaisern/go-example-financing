package esrcwatermill

import (
	"context"

	"github.com/pperaltaisern/financing/internal/esrc/relay"

	"github.com/ThreeDotsLabs/watermill/components/cqrs"
)

type Publisher struct {
	bus *cqrs.EventBus
}

var _ relay.Publisher = (*Publisher)(nil)

func NewPublisher(bus *cqrs.EventBus) Publisher {
	return Publisher{bus}
}

func (p Publisher) Publish(ctx context.Context, e relay.RelayEvent) error {
	return p.bus.Publish(ctx, e)
}
