package e2e

import (
	"testing"

	"github.com/pperaltaisern/financing/pkg/financing"
	"github.com/pperaltaisern/financing/pkg/intevent"
)

func (s *CommandsSuite) TestIssuerRegistered() {
	t := s.T()

	t.Run("GIVEN a IssuerRegistered integration event WHEN we process it THEN we create the issuer", func(t *testing.T) {
		issuerID := financing.NewID()
		s.RegisterIssuer(issuerID)
	})
}

func (s *CommandsSuite) RegisterIssuer(id financing.ID) {
	t := s.T()

	integrationEvent := intevent.IssuerRegistered{
		ID:   id,
		Name: "ISSUER_1",
	}
	eventAssertion := EventAssertion{
		Expected: financing.NewIssuerCreatedEvent(id),
		Actual:   &financing.IssuerCreatedEvent{},
	}

	s.publishIntegrationEventAndAssertCreatedInEventSource(t, id, integrationEvent)

	s.expectEvents(t, eventAssertion)
}
