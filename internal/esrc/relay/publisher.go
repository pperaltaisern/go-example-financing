package relay

import "context"

type Publisher interface {
	Publish(context.Context, RelayEvent) error
}
