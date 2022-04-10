package acceptance

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/pperaltaisern/esrc"
	"github.com/pperaltaisern/esrcwatermill"
	"github.com/pperaltaisern/financing/pkg/config"
	"github.com/pperaltaisern/financing/pkg/financing"
	"github.com/pperaltaisern/financing/pkg/grpc/pb"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type CommandsSuite struct {
	suite.Suite

	conn               *grpc.ClientConn
	commands           pb.CommandsClient
	eventStore         esrc.EventStore
	eventBus           *cqrs.EventBus
	subscriberMessageC <-chan *message.Message
	waitTime           time.Duration
	eventMarshaler     cqrs.CommandEventMarshaler
}

func TestCommandFeatures(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping acceptance test")
	}

	s := new(CommandsSuite)
	suite.Run(t, s)
}

func (s *CommandsSuite) SetupSuite() {
	s.waitTime = time.Second

	log, err := config.LoadLoggerConfig().Build()
	if err != nil {
		panic(err)
	}

	serverConfig := config.LoadCommandServerConfig()
	conn, err := grpc.Dial(serverConfig.Address, grpc.WithInsecure())
	if err != nil {
		log.Fatal("err building commands client", zap.Error(err))
	}
	s.conn = conn
	s.commands = pb.NewCommandsClient(conn)

	repos, es, err := config.LoadCommandPostgresConfig().BuildRepositories()
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

	s.purgeMessageQueue()
}

func (s *CommandsSuite) TearDownTest() {
	s.AssertNoMoreMessages(s.T())
}

func (s *CommandsSuite) TearDownSuite() {
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
func (s *CommandsSuite) expectEvents(t *testing.T, events ...EventAssertion) {
	for _, e := range events {
		var m *message.Message
		t.Logf("expecting event '%s'", e.Actual.EventName())
		for {
			m = s.waitForMessage(t)
			eventName := s.eventMarshaler.NameFromMessage(m)
			// integrationEvents should come from another queue but for simplicity we put them in the same and ignore them when asserting
			if _, ok := integrationEvents[eventName]; !ok {
				break
			}
		}
		require.Equal(t, e.Actual.EventName(), s.eventMarshaler.NameFromMessage(m))

		err := s.eventMarshaler.Unmarshal(m, e.Actual)
		require.NoError(t, err)
		require.Equal(t, e.Expected, e.Actual)
	}
}

func (s *CommandsSuite) waitForMessage(t *testing.T) *message.Message {
	for i := 0; i < 3; i++ {
		select {
		case m := <-s.subscriberMessageC:
			s.T().Logf("message obtained from queue: %s", s.eventMarshaler.NameFromMessage(m))
			m.Ack()
			return m
		default:
			t.Log("wait for message...")
			time.Sleep(s.waitTime)
		}
	}
	require.FailNow(t, "message not received")
	return nil
}

func (s *CommandsSuite) AssertNoMoreMessages(t *testing.T) {
	time.Sleep(s.waitTime)
	var messages []*message.Message
Loop:
	for {
		select {
		case m := <-s.subscriberMessageC:
			s.T().Logf("message obtained from queue in teardown: %s", s.eventMarshaler.NameFromMessage(m))
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

func (s *CommandsSuite) purgeMessageQueue() {
	for {
		select {
		case m := <-s.subscriberMessageC:
			s.T().Logf("message purged: %s", s.eventMarshaler.NameFromMessage(m))
			m.Ack()
			time.Sleep(s.waitTime)
		default:
			return
		}
	}
}

func (s *CommandsSuite) publishIntegrationEventAndAssertCreatedInEventSource(t *testing.T, id financing.ID, event interface{}) {
	s.publishEvent(t, event)
	s.assertContains(t, id)
}

func (s *CommandsSuite) publishEvent(t *testing.T, event interface{}) {
	err := s.eventBus.Publish(context.Background(), event)
	require.NoError(t, err)
	// wait the event to be processed
	time.Sleep(s.waitTime)
}

func (s *CommandsSuite) assertContains(t *testing.T, id financing.ID) {
	s.assertContainsBool(t, id, true)
}

func (s *CommandsSuite) assertNotContains(t *testing.T, id financing.ID) {
	s.assertContainsBool(t, id, false)
}

func (s *CommandsSuite) assertContainsBool(t *testing.T, id financing.ID, expected bool) {
	// wait the event to be processed
	time.Sleep(s.waitTime)
	f, err := s.eventStore.ContainsAggregate(context.Background(), "", id)
	require.NoError(t, err)
	require.Equal(t, expected, f)
}
