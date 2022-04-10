package projection

import (
	"testing"

	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/pperaltaisern/esrcwatermill"
	"github.com/pperaltaisern/financing/pkg/financing"
	"github.com/stretchr/testify/require"
)

type UnknownEvent struct{}

func (UnknownEvent) EventName() string { return "UnknownEvent" }

func TestProjectMessageShould(t *testing.T) {
	eventMarshaller := esrcwatermill.RelayEventMarshaler{
		CmdMarshaler: cqrs.JSONMarshaler{
			GenerateName: cqrs.StructName,
		},
	}

	knownEvent := financing.NewInvestorCreatedEvent(financing.NewID())
	unknownEvent := UnknownEvent{}

	t.Run("CallProjectEventAndACKForKnownEvent", func(t *testing.T) {
		// given
		var capturedEvent *financing.InvestorCreatedEvent
		eventProjector := MockEventProjector{
			ProjectInvestorCreatedEventFn: func(arg *financing.InvestorCreatedEvent) error {
				capturedEvent = arg
				return nil
			},
		}

		knownMessage, err := eventMarshaller.Marshal(knownEvent)
		require.NoError(t, err)

		messageProjector := NewMessageProjector(&eventProjector, eventMarshaller, nil)

		// when
		messageProjector.ProjectMessage(knownMessage)

		// then
		require.Equal(t, knownEvent, capturedEvent)
		ackedC := knownMessage.Acked()
		_, isOpen := <-ackedC
		require.False(t, isOpen)
	})

	t.Run("LogErrorAndAckForUnknownEvent", func(t *testing.T) {
		// given
		eventProjector := MockEventProjector{}

		logged := false
		logErr := func(m *message.Message, e error) { logged = true }

		unknownMessage, err := eventMarshaller.Marshal(unknownEvent)
		require.NoError(t, err)

		messageProjector := NewMessageProjector(&eventProjector, eventMarshaller, logErr)

		// when
		messageProjector.ProjectMessage(unknownMessage)

		// then
		require.True(t, logged)
		ackedC := unknownMessage.Acked()
		_, isOpen := <-ackedC
		require.False(t, isOpen)
	})
}
