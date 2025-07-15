package order

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	clientMocks "github.com/Alexey-step/rocket-factory/order/internal/client/grpc/mocks"
	"github.com/Alexey-step/rocket-factory/order/internal/repository/mocks"
)

type ServiceSuite struct {
	suite.Suite
	//nolint:containedctx
	ctx context.Context

	orderRepository *mocks.OrderRepository
	inventoryClient *clientMocks.InventoryClient
	paymentClient   *clientMocks.PaymentClient

	service *service
}

func (s *ServiceSuite) SetupTest() {
	s.ctx = context.Background()

	s.orderRepository = mocks.NewOrderRepository(s.T())
	s.inventoryClient = clientMocks.NewInventoryClient(s.T())
	s.paymentClient = clientMocks.NewPaymentClient(s.T())

	s.service = NewService(
		s.orderRepository,
		s.inventoryClient,
		s.paymentClient,
	)
}

func (s *ServiceSuite) TearDownTest() {
}

func TestServiceIntegration(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}
