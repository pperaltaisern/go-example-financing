package esrc

type RawEvent struct {
	Name string
	Data []byte
}

type Event interface {
	EventName() string
}
