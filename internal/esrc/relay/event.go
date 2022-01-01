package relay

type Event struct {
	ID   uint64
	Name string
	Data []byte
}

func NewEvent(id uint64, name string, data []byte) Event {
	return Event{
		ID:   id,
		Name: name,
		Data: data,
	}
}
