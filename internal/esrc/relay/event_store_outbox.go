package relay

import "context"

type EventStoreOutbox interface {
	UnpublishedEvents(context.Context) ([]RelayEvent, error)
	MarkEventsAsPublised(context.Context, []RelayEvent) error
}
