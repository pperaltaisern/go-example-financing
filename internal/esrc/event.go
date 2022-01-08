package esrc

import "encoding/json"

type RawEvent struct {
	Name string
	Data []byte
}

type Event interface {
	EventName() string
}

type MarshalEvent func(Event) ([]byte, error)

// MarshalEvents marshals Events to RawEvents (in the same order) given a MarshalEvent func (as could be json.Marshal),
// if an error happens marshalling an Event, the events marshalled succesfully before that event are returned along with the error.
func MarshalEvents(events []Event, marshal MarshalEvent) ([]RawEvent, error) {
	rawEvents := make([]RawEvent, 0, len(events))
	for _, e := range events {
		b, err := marshal(e)
		if err != nil {
			return rawEvents, err
		}
		rawEvents = append(rawEvents, RawEvent{Name: e.EventName(), Data: b})
	}
	return rawEvents, nil
}

func MarshalEventsJSON(events []Event) ([]RawEvent, error) {
	return MarshalEvents(events, marshalEventJson)
}

func marshalEventJson(e Event) ([]byte, error) {
	return json.Marshal(e)
}
