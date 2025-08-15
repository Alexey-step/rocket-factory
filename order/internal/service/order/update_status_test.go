package order

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"

	clientMocks "github.com/Alexey-step/rocket-factory/order/internal/client/grpc/mocks"
	"github.com/Alexey-step/rocket-factory/order/internal/model"
	"github.com/Alexey-step/rocket-factory/order/internal/repository/mocks"
	orderServiceMocks "github.com/Alexey-step/rocket-factory/order/internal/service/mocks"
	"github.com/Alexey-step/rocket-factory/platform/pkg/logger"
)

func TestUpdateStatusSuccess(t *testing.T) {
	ctx := context.Background()
	orderUUID := gofakeit.UUID()
	status := model.OrderStatusCompleted

	orderRepository := mocks.NewOrderRepository(t)
	inventoryClient := clientMocks.NewInventoryClient(t)
	paymentClient := clientMocks.NewPaymentClient(t)
	orderProducer := orderServiceMocks.NewOrderProducerService(t)

	orderService := NewService(
		orderRepository,
		inventoryClient,
		paymentClient,
		orderProducer,
	)

	orderInfo := model.OrderUpdateInfo{
		Status: lo.ToPtr(model.OrderStatusCompleted),
	}

	orderRepository.On("UpdateOrder", ctx, orderUUID, orderInfo).Return(nil).Once()

	err := orderService.UpdateStatus(ctx, orderUUID, status)
	assert.NoError(t, err)
}

func TestUpdateStatusFail(t *testing.T) {
	ctx := context.Background()
	orderUUID := gofakeit.UUID()
	status := model.OrderStatusCompleted
	expectedErr := gofakeit.Error()

	logger.SetNopLogger()
	orderRepository := mocks.NewOrderRepository(t)
	inventoryClient := clientMocks.NewInventoryClient(t)
	paymentClient := clientMocks.NewPaymentClient(t)
	orderProducer := orderServiceMocks.NewOrderProducerService(t)

	orderService := NewService(
		orderRepository,
		inventoryClient,
		paymentClient,
		orderProducer,
	)

	orderInfo := model.OrderUpdateInfo{
		Status: lo.ToPtr(model.OrderStatusCompleted),
	}

	orderRepository.On("UpdateOrder", ctx, orderUUID, orderInfo).Return(expectedErr).Once()

	err := orderService.UpdateStatus(ctx, orderUUID, status)
	assert.Error(t, err)
	assert.Equal(t, err, expectedErr)
}
