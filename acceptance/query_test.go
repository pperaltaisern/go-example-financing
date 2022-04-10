package acceptance

import (
	"context"
	"testing"
	"time"

	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/pperaltaisern/esrc"
	"github.com/pperaltaisern/esrc/relay"
	"github.com/pperaltaisern/financing/pkg/config"
	"github.com/pperaltaisern/financing/pkg/financing"
	"github.com/pperaltaisern/financing/pkg/grpc/pb"
	"github.com/pperaltaisern/financing/pkg/query/pg"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
)

type QueriesSuite struct {
	suite.Suite

	conn           *grpc.ClientConn
	queries        pb.QueriesClient
	eventBus       *cqrs.EventBus
	waitTime       time.Duration
	eventMarshaler cqrs.CommandEventMarshaler
	eventProjector *pg.EventProjector
}

func TestQueryFeatures(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping acceptance test")
	}

	s := new(QueriesSuite)
	suite.Run(t, s)
}

func (s *QueriesSuite) SetupSuite() {
	s.waitTime = time.Second

	log, err := config.LoadLoggerConfig().Build()
	require.NoError(s.T(), err)

	serverConfig := config.LoadQueryServerConfig()
	conn, err := grpc.Dial(serverConfig.Address, grpc.WithInsecure())
	require.NoError(s.T(), err)

	s.conn = conn
	s.queries = pb.NewQueriesClient(conn)

	eventBus, err := config.LoadAMQPConfig().BuildEventBus(log)
	require.NoError(s.T(), err)

	s.eventBus = eventBus
	s.eventMarshaler = config.CommandEventMarshaler
}

func (s *QueriesSuite) SetupTest() {
	db, err := config.LoadQueryPostgresConfig().BuildGORM()
	require.NoError(s.T(), err)

	s.eventProjector, err = pg.NewEventProjector(db)
	require.NoError(s.T(), err)

	err = s.eventProjector.Clean()
	require.NoError(s.T(), err)
}

func (s *QueriesSuite) TearDownTest() {
	err := s.eventProjector.Clean()
	require.NoError(s.T(), err)
}

func (s *QueriesSuite) TearDownSuite() {
	s.conn.Close()
}

func (s *QueriesSuite) newRelayEvent(aggregateID financing.ID, event esrc.Event) relay.RelayEvent {
	rawEvents, err := esrc.MarshalEvents([]esrc.Event{event}, esrc.JSONEventMarshaler{})
	require.NoError(s.T(), err)
	return relay.NewRelayEvent(aggregateID, 0, rawEvents[0])
}

func (s *QueriesSuite) publisEvents(events ...relay.RelayEvent) {
	for _, e := range events {
		err := s.eventBus.Publish(context.Background(), e)
		require.NoError(s.T(), err)
	}
	time.Sleep(s.waitTime)
}
