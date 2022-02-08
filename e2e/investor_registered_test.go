package e2e

import (
	"testing"

	"github.com/pperaltaisern/financing/pkg/financing"
	"github.com/pperaltaisern/financing/pkg/intevent"
)

func (s *Suite) TestInvestorRegistered() {
	t := s.T()

	t.Run("GIVEN a InvestorRegistered integration event WHEN we process it THEN we create the investor", func(t *testing.T) {
		investorID := financing.NewID()
		s.RegisterInvestor(investorID, 100)
	})
}

func (s *Suite) RegisterInvestor(id financing.ID, balance financing.Money) {
	t := s.T()

	integrationEvent := intevent.InvestorRegistered{
		ID:      id,
		Name:    "INVESTOR_1",
		Balance: 100,
	}
	eventAssertions := []EventAssertion{
		{
			Expected: financing.NewInvestorCreatedEvent(id),
			Actual:   &financing.InvestorCreatedEvent{},
		},
		{
			Expected: financing.NewInvestorFundsAddedEvent(id, balance),
			Actual:   &financing.InvestorFundsAddedEvent{},
		},
	}

	s.publishIntegrationEventAndAssertCreatedInEventSource(t, id, integrationEvent)

	s.expectEvents(t, eventAssertions...)
}
