package financing

type IssuerCreatedEvent struct {
	IssuerID ID
}

func NewIssuerCreatedEvent(issuerID ID) IssuerCreatedEvent {
	return IssuerCreatedEvent{
		IssuerID: issuerID,
	}
}

func (e IssuerCreatedEvent) IsEvent() {}

func (e IssuerCreatedEvent) Name() string { return "IssuerCreatedEvent" }
