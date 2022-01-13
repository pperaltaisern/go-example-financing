package relay

import "github.com/pperaltaisern/financing/internal/esrc"

type RelayEvent struct {
	AggregateID esrc.ID
	Sequence    uint64
	RawEvent    esrc.RawEvent
}

func NewRelayEvent(aggregateID esrc.ID, sequence uint64, e esrc.RawEvent) RelayEvent {
	return RelayEvent{
		AggregateID: aggregateID,
		Sequence:    sequence,
		RawEvent:    e,
	}
}
