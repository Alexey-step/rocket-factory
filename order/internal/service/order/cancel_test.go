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
)

func TestCancelOrderSuccess(t *testing.T) {
	ctx := context.Background()
	orderUUID := gofakeit.UUID()
	order := getMockOrderWithStatus(orderUUID, model.OrderStatusPendingPayment)

	orderUpdateInfo := model.OrderUpdateInfo{
		Status: lo.ToPtr(model.OrderStatusCanceled),
	}

	orderRepository := mocks.NewOrderRepository(t)
	inventoryClient := clientMocks.NewInventoryClient(t)
	paymentClient := clientMocks.NewPaymentClient(t)

	orderService := NewService(
		orderRepository,
		inventoryClient,
		paymentClient,
	)

	orderRepository.On("GetOrder", ctx, orderUUID).Return(order, nil).Once()
	orderRepository.On("UpdateOrder", ctx, orderUUID, orderUpdateInfo).Return(nil).Once()
	err := orderService.CancelOrder(ctx, orderUUID)
	assert.NoError(t, err)
}

func TestCancelOrderFail(t *testing.T) {
	ctx := context.Background()
	orderUUID := gofakeit.UUID()
	expectedErr := model.ErrOrderConflict

	order := getMockOrderWithStatus(orderUUID, model.OrderStatusPaid)

	orderRepository := mocks.NewOrderRepository(t)
	inventoryClient := clientMocks.NewInventoryClient(t)
	paymentClient := clientMocks.NewPaymentClient(t)

	orderService := NewService(
		orderRepository,
		inventoryClient,
		paymentClient,
	)

	orderRepository.On("GetOrder", ctx, orderUUID).Return(order, nil).Once()
	err := orderService.CancelOrder(ctx, orderUUID)
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
}

func TestCancelOrderConflictFail(t *testing.T) {
	ctx := context.Background()
	orderUUID := gofakeit.UUID()
	expectedErr := model.ErrOrderConflict

	order := getMockOrderWithStatus(orderUUID, model.OrderStatusCanceled)

	orderRepository := mocks.NewOrderRepository(t)
	inventoryClient := clientMocks.NewInventoryClient(t)
	paymentClient := clientMocks.NewPaymentClient(t)

	orderService := NewService(
		orderRepository,
		inventoryClient,
		paymentClient,
	)

	orderRepository.On("GetOrder", ctx, orderUUID).Return(order, nil).Once()
	err := orderService.CancelOrder(ctx, orderUUID)
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
}

func TestCancelOrderInternalError(t *testing.T) {
	ctx := context.Background()
	orderUUID := gofakeit.UUID()
	expectedErr := model.ErrOrderInternalError

	orderRepository := mocks.NewOrderRepository(t)
	inventoryClient := clientMocks.NewInventoryClient(t)
	paymentClient := clientMocks.NewPaymentClient(t)

	orderService := NewService(
		orderRepository,
		inventoryClient,
		paymentClient,
	)

	orderRepository.On("GetOrder", ctx, orderUUID).Return(model.OrderData{}, expectedErr).Once()
	err := orderService.CancelOrder(ctx, orderUUID)
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
}

func TestCancelOrderNotFoundFail(t *testing.T) {
	ctx := context.Background()
	orderUUID := gofakeit.UUID()
	expectedErr := model.ErrOrderNotFound

	orderRepository := mocks.NewOrderRepository(t)
	inventoryClient := clientMocks.NewInventoryClient(t)
	paymentClient := clientMocks.NewPaymentClient(t)

	orderService := NewService(
		orderRepository,
		inventoryClient,
		paymentClient,
	)

	orderRepository.On("GetOrder", ctx, orderUUID).Return(model.OrderData{}, expectedErr).Once()
	err := orderService.CancelOrder(ctx, orderUUID)
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
}

func TestCancelOrderUpdateError(t *testing.T) {
	ctx := context.Background()
	orderUUID := gofakeit.UUID()
	expectedErr := model.ErrOrderNotFound

	order := getMockOrderWithStatus(orderUUID, model.OrderStatusPendingPayment)

	orderUpdateInfo := model.OrderUpdateInfo{
		Status: lo.ToPtr(model.OrderStatusCanceled),
	}

	orderRepository := mocks.NewOrderRepository(t)
	inventoryClient := clientMocks.NewInventoryClient(t)
	paymentClient := clientMocks.NewPaymentClient(t)

	orderService := NewService(
		orderRepository,
		inventoryClient,
		paymentClient,
	)

	orderRepository.On("GetOrder", ctx, orderUUID).Return(order, nil).Once()
	orderRepository.On("UpdateOrder", ctx, orderUUID, orderUpdateInfo).Return(expectedErr).Once()
	err := orderService.CancelOrder(ctx, orderUUID)
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
}

func TestCancelOrderInternalErr(t *testing.T) {
	ctx := context.Background()
	orderUUID := gofakeit.UUID()
	expectedErr := model.ErrOrderInternalError

	order := getMockOrderWithStatus(orderUUID, "UNKNOWN_STATUS")

	orderRepository := mocks.NewOrderRepository(t)
	inventoryClient := clientMocks.NewInventoryClient(t)
	paymentClient := clientMocks.NewPaymentClient(t)

	orderService := NewService(
		orderRepository,
		inventoryClient,
		paymentClient,
	)

	orderRepository.On("GetOrder", ctx, orderUUID).Return(order, nil).Once()
	err := orderService.CancelOrder(ctx, orderUUID)
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
}

func getMockOrderWithStatus(orderUUID string, status model.OrderStatus) model.OrderData {
	return model.OrderData{
		UUID:          orderUUID,
		UserUUID:      gofakeit.UUID(),
		PartUuids:     []string{gofakeit.UUID()},
		TotalPrice:    gofakeit.Price(100, 1000),
		PaymentMethod: lo.ToPtr(model.PaymentMethod("CREDIT_CARD")),
		Status:        status,
		CreatedAt:     gofakeit.Date(),
	}
}
