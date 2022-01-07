package financing

type IssuerCreatedEvent struct {
	IssuerID ID `json:"-"`
}

func NewIssuerCreatedEvent(issuerID ID) *IssuerCreatedEvent {
	return &IssuerCreatedEvent{
		IssuerID: issuerID,
	}
}

func (e *IssuerCreatedEvent) EventName() string { return "IssuerCreatedEvent" }

func (e *IssuerCreatedEvent) WithAggregateID(id string) {
	e.IssuerID = NewIDFrom(id)
}
