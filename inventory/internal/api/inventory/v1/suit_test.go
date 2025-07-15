package v1

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/Alexey-step/rocket-factory/inventory/internal/service/mocks"
)

type ServiceSuite struct {
	suite.Suite

	ctx context.Context

	service *mocks.InventoryService
	api     *api
}

func (s *ServiceSuite) SetupTest() {
	s.ctx = context.Background()

	s.service = mocks.NewInventoryService(s.T())
	s.api = NewAPI(s.service)
}

func (s *ServiceSuite) TearDownTest() {
}

func TestServiceIntegration(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}
