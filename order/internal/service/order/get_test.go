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

func TestGetOrderSuccess(t *testing.T) {
	ctx := context.Background()
	orderUUID := gofakeit.UUID()

	order := getMockedOrder(orderUUID)

	orderRepository := mocks.NewOrderRepository(t)
	inventoryClient := clientMocks.NewInventoryClient(t)
	paymentClient := clientMocks.NewPaymentClient(t)

	orderService := NewService(
		orderRepository,
		inventoryClient,
		paymentClient,
	)

	orderRepository.On("GetOrder", ctx, orderUUID).Return(order, nil).Once()
	resp, err := orderService.GetOrder(ctx, orderUUID)

	assert.NoError(t, err)
	assert.Equal(t, order, resp)
}

func TestGetOrderNotFoundErr(t *testing.T) {
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
	resp, err := orderService.GetOrder(ctx, orderUUID)

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.Empty(t, resp)
}

func TestGetOrderInternalErr(t *testing.T) {
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
	resp, err := orderService.GetOrder(ctx, orderUUID)

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.Empty(t, resp)
}

func getMockedOrder(uuid string) model.OrderData {
	return model.OrderData{
		UUID:          uuid,
		UserUUID:      gofakeit.UUID(),
		PartUuids:     []string{gofakeit.UUID()},
		TotalPrice:    gofakeit.Price(100, 1000),
		PaymentMethod: lo.ToPtr(model.PaymentMethod("CREDIT_CARD")),
		Status:        model.OrderStatusPendingPayment,
		CreatedAt:     gofakeit.Date(),
	}
}
