package e2e

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/pperaltaisern/financing/internal/esrc"
	"github.com/pperaltaisern/financing/internal/esrc/esrcwatermill"
	"github.com/pperaltaisern/financing/pkg/config"
	"github.com/pperaltaisern/financing/pkg/financing"
	"github.com/pperaltaisern/financing/pkg/grpc/pb"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type Suite struct {
	suite.Suite

	conn               *grpc.ClientConn
	commands           pb.CommandsClient
	eventStore         esrc.EventStore
	eventBus           *cqrs.EventBus
	subscriberMessageC <-chan *message.Message
	waitTime           time.Duration
	eventMarshaler     cqrs.CommandEventMarshaler
}

func TestFeatures(t *testing.T) {
	s := new(Suite)
	suite.Run(t, s)
}

func (s *Suite) SetupSuite() {
	s.waitTime = time.Second

	log, err := config.LoadLoggerConfig().Build()
	if err != nil {
		panic(err)
	}

	serverConfig := config.LoadServerConfig()
	conn, err := grpc.Dial(serverConfig.Address, grpc.WithInsecure())
	if err != nil {
		log.Fatal("err building commands client", zap.Error(err))
	}
	s.conn = conn
	s.commands = pb.NewCommandsClient(conn)

	repos, es, err := config.LoadPostgresConfig().BuildRepositories()
	if err != nil {
		log.Fatal("err building repositories", zap.Error(err))
	}
	s.eventStore = es

	facade, _, err := config.BuildCqrsFacade(log, repos)
	if err != nil {
		log.Fatal("err building facade", zap.Error(err))
	}
	s.eventBus = facade.EventBus()

	const testEventSubscriberTopic = "test_handler"
	eventSubscriber, err := config.LoadAMQPConfig().BuildEventSubscriber(log, testEventSubscriberTopic)
	if err != nil {
		log.Fatal("err building test subscriber", zap.Error(err))
	}
	s.subscriberMessageC, err = eventSubscriber.Subscribe(context.Background(), "events")
	if err != nil {
		log.Fatal("err subscribing to test subscriber", zap.Error(err))
	}

	s.eventMarshaler = esrcwatermill.RelayEventMarshaler{
		CmdMarshaler: cqrs.JSONMarshaler{},
	}
}

func (s *Suite) TearDownSuite() {
	s.AssertNoMoreMessages(s.T())
	s.conn.Close()
}

type EventAssertion struct {
	Expected esrc.Event
	// Actual is an empty event that is going to be used to unmarshal the message
	Actual esrc.Event
}

var integrationEvents = map[string]struct{}{
	"InvestorRegistered": {},
	"IssuerRegistered":   {},
}

// expectEvents asserts that messages received from the subscribed are the expected, the order is important
func (s *Suite) expectEvents(t *testing.T, events ...EventAssertion) {
	for _, e := range events {
		var m *message.Message
		for {
			m = s.waitForMessage(t)
			eventName := m.Metadata.Get("name")
			// integrationEvents should come from another queue but for simplicity we put them in the same and ignore them when asserting
			if _, ok := integrationEvents[eventName]; !ok {
				break
			}
		}
		require.Equal(t, e.Actual.EventName(), m.Metadata.Get("name"))

		err := s.eventMarshaler.Unmarshal(m, e.Actual)
		require.NoError(t, err)
		require.Equal(t, e.Expected, e.Actual)
	}
}

func (s *Suite) waitForMessage(t *testing.T) *message.Message {
	for i := 0; i < 3; i++ {
		select {
		case m := <-s.subscriberMessageC:
			m.Ack()
			return m
		default:
			time.Sleep(s.waitTime)
		}
	}
	require.FailNow(t, "message not received")
	return nil
}

func (s *Suite) AssertNoMoreMessages(t *testing.T) {
	time.Sleep(s.waitTime)
	var messages []*message.Message
Loop:
	for {
		select {
		case m := <-s.subscriberMessageC:
			m.Ack()
			messages = append(messages, m)
		default:
			break Loop
		}
	}
	if len(messages) > 0 {
		b, err := json.Marshal(messages)
		require.FailNowf(t, "shouldn't be more messages in queue after all tests are finished, found:", string(b), err)
	}
}

func (s *Suite) publishIntegrationEventAndAssertCreatedInEventSource(t *testing.T, id financing.ID, event interface{}) {
	s.publishEvent(t, event)
	s.assertContains(t, id)
}

func (s *Suite) publishEvent(t *testing.T, event interface{}) {
	err := s.eventBus.Publish(context.Background(), event)
	require.NoError(t, err)
	// wait the event to be processed
	time.Sleep(s.waitTime)
}

func (s *Suite) assertContains(t *testing.T, id financing.ID) {
	s.assertContainsBool(t, id, true)
}

func (s *Suite) assertNotContains(t *testing.T, id financing.ID) {
	s.assertContainsBool(t, id, false)
}

func (s *Suite) assertContainsBool(t *testing.T, id financing.ID, expected bool) {
	// wait the event to be processed
	time.Sleep(s.waitTime)
	f, err := s.eventStore.Contains(context.Background(), id)
	require.NoError(t, err)
	require.Equal(t, expected, f)
}
