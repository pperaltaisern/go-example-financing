package relay

import "context"

type EventStoreOutbox interface {
	UnpublishedEvents(context.Context) ([]Event, error)
	MarkEventsAsPublised(context.Context, []Event) error
}
