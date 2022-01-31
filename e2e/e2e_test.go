package e2e

import (
	"context"
	"testing"
	"time"

	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/pperaltaisern/financing/internal/esrc"
	"github.com/pperaltaisern/financing/pkg/config"
	"github.com/pperaltaisern/financing/pkg/financing"
	"github.com/pperaltaisern/financing/pkg/grpc/pb"
	"github.com/pperaltaisern/financing/pkg/intevent"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type Suite struct {
	suite.Suite

	conn       *grpc.ClientConn
	commands   pb.CommandsClient
	eventStore esrc.EventStore
	eventBus   *cqrs.EventBus
	waitTime   time.Duration

	investorID financing.ID
	issuerID   financing.ID
	invoiceID  financing.ID
}

func (s *Suite) SetupTest() {
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

	s.investorID = financing.NewID()
	s.issuerID = financing.NewIDFromString("37bca316-3b73-4caf-8230-6e4f287ab2e1")
	s.investorID = financing.NewID()
}

func (s *Suite) TearDownTest() {
	s.conn.Close()
}

func (s *Suite) TestInvestorRegistered() {
	t := s.T()

	e := intevent.InvestorRegistered{
		ID:      s.investorID,
		Name:    "INVESTOR_1",
		Balance: 100,
	}
	s.publishIntegrationEvent(t, s.investorID, e)
}

func (s *Suite) TestIssuerRegistered() {
	t := s.T()

	e := intevent.IssuerRegistered{
		ID:   s.issuerID,
		Name: "ISSUER_1",
	}
	s.publishIntegrationEvent(t, s.issuerID, e)
}

func (s *Suite) TestSellInvoice() {
	t := s.T()

	cmd := &pb.SellInvoiceCommand{
		IssuerId: &pb.UUID{
			Value: s.issuerID.String(),
		},
		AskingPrice: &pb.Money{
			Amount: 20,
		},
	}

	id, err := s.commands.SellInvoice(context.Background(), cmd)
	require.NoError(t, err)

	s.assertContains(t, s.invoiceID)

	s.invoiceID = financing.NewIDFromString(id.Value)
}

// func (s *Suite) TestSellInvoice_NotExistingIssuer() {
// 	t := s.T()

// 	cmd := &pb.SellInvoiceCommand{
// 		IssuerId: &pb.UUID{
// 			// this issuer won't be found
// 			Value: financing.NewID().String(),
// 		},
// 		AskingPrice: &pb.Money{
// 			Amount: 20,
// 		},
// 	}

// 	id, err := s.commands.SellInvoice(context.Background(), cmd)
// 	require.NoError(t, err)

// 	s.assertNotContains(t, financing.NewIDFromString(id.Value))
// }

// func (s *Suite) TestBidOnInvoice() {
// 	t := s.T()

// 	cmd := &pb.BidOnInvoiceCommand{
// 		InvestorId: &pb.UUID{
// 			Value: s.investorID.String(),
// 		},
// 		InvoiceId: &pb.UUID{
// 			Value: s.invoiceID.String(),
// 		},
// 		Bid: &pb.Money{
// 			Amount: 30,
// 		},
// 	}

// 	_, err := s.commands.BidOnInvoice(context.Background(), cmd)
// 	require.NoError(t, err)

// 	s.eventBus()
// }

func (s *Suite) publishIntegrationEvent(t *testing.T, id financing.ID, event interface{}) {
	err := s.eventBus.Publish(context.Background(), event)
	require.NoError(t, err)

	s.assertContains(t, id)
}

func (s *Suite) assertContains(t *testing.T, id financing.ID) {
	s.assertContainsBool(t, id, true)
}

func (s *Suite) assertNotContains(t *testing.T, id financing.ID) {
	s.assertContainsBool(t, id, false)
}

func (s *Suite) assertContainsBool(t *testing.T, id financing.ID, expected bool) {
	time.Sleep(s.waitTime)

	f, err := s.eventStore.Contains(context.Background(), id)
	require.NoError(t, err)
	require.Equal(t, expected, f)
}

func TestFeatures(t *testing.T) {
	suite.Run(t, new(Suite))
}
