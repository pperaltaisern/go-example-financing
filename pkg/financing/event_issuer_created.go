package financing

type IssuerCreatedEvent struct {
	IssuerID ID
}

func NewIssuerCreatedEvent(issuerID ID) *IssuerCreatedEvent {
	return &IssuerCreatedEvent{
		IssuerID: issuerID,
	}
}

func (e *IssuerCreatedEvent) EventName() string { return "IssuerCreatedEvent" }
