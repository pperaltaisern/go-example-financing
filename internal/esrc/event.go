package esrc

import "encoding/json"

// Event must be implemented by domain events
type Event interface {
	EventName() string
}

// RawEvent is the marshaled version of an Event
type RawEvent struct {
	Name string
	Data []byte
}

// EventFactory is the abstraction needed to create the Event associated to an EventName
type EventFactory interface {
	CreateEmptyEvent(name string) (Event, error)
}

// EventMarshaler marshals and unmarshals an Event's data/body
type EventMarshaler interface {
	MarshalEvent(Event) ([]byte, error)
	UnmarshalEvent([]byte, Event) error
}

type JSONEventMarshaler struct{}

func (JSONEventMarshaler) MarshalEvent(e Event) ([]byte, error) {
	return json.Marshal(e)
}

func (JSONEventMarshaler) UnmarshalEvent(b []byte, e Event) error {
	return json.Unmarshal(b, e)
}

// MarshalEvents marshals Events to RawEvents (in the same order) given an EventMarshaler,
// if an error happens marshalling an Event, the events marshalled succesfully before that event are returned along with the error.
func MarshalEvents(events []Event, marshaler EventMarshaler) ([]RawEvent, error) {
	rawEvents := make([]RawEvent, 0, len(events))
	for _, e := range events {
		b, err := marshaler.MarshalEvent(e)
		if err != nil {
			return rawEvents, err
		}
		rawEvents = append(rawEvents, RawEvent{Name: e.EventName(), Data: b})
	}
	return rawEvents, nil
}
