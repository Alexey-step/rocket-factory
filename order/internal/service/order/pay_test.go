package order

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	clientMocks "github.com/Alexey-step/rocket-factory/order/internal/client/grpc/mocks"
	"github.com/Alexey-step/rocket-factory/order/internal/model"
	"github.com/Alexey-step/rocket-factory/order/internal/repository/mocks"
	orderServiceMocks "github.com/Alexey-step/rocket-factory/order/internal/service/mocks"
)

func TestPayOrderSuccess(t *testing.T) {
	ctx := context.Background()
	orderUUID := gofakeit.UUID()
	paymentMethod := "CREDIT_CARD"
	transactionUUID := gofakeit.UUID()

	order := getMockedPayOrder(orderUUID, model.OrderStatusPendingPayment)

	orderInfo := model.OrderUpdateInfo{
		Status:          lo.ToPtr(model.OrderStatusPaid),
		PaymentMethod:   lo.ToPtr(model.PaymentMethod(paymentMethod)),
		TransactionUUID: lo.ToPtr(transactionUUID),
	}

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

	orderRepository.On("GetOrder", mock.Anything, orderUUID).Return(order, nil).Once()
	paymentClient.On("PayOrder", mock.Anything, order.UserUUID, orderUUID, paymentMethod).Return(transactionUUID, nil).Once()
	orderProducer.On("ProduceOrderPaid", mock.Anything, mock.AnythingOfType("model.OrderPaid")).Return(nil).Once()
	orderRepository.On("UpdateOrder", mock.Anything, orderUUID, orderInfo).Return(nil).Once()

	resp, err := orderService.PayOrder(ctx, orderUUID, paymentMethod)
	assert.NoError(t, err)
	assert.Equal(t, transactionUUID, resp)
}

func TestPayOrderFailGetOrder(t *testing.T) {
	ctx := context.Background()
	orderUUID := gofakeit.UUID()
	paymentMethod := "CREDIT_CARD"
	expectedErr := model.ErrOrderNotFound

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

	orderRepository.On("GetOrder", mock.Anything, orderUUID).Return(model.OrderData{}, expectedErr).Once()

	resp, err := orderService.PayOrder(ctx, orderUUID, paymentMethod)
	assert.Error(t, err)
	assert.Equal(t, err, expectedErr)
	assert.Empty(t, resp)
}

func TestPayOrderFail(t *testing.T) {
	ctx := context.Background()
	orderUUID := gofakeit.UUID()
	paymentMethod := "CREDIT_CARD"
	expectedErr := model.ErrPaymentConflict

	order := getMockedPayOrder(orderUUID, model.OrderStatusPendingPayment)

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

	orderRepository.On("GetOrder", mock.Anything, orderUUID).Return(order, nil).Once()
	paymentClient.On("PayOrder", mock.Anything, order.UserUUID, orderUUID, paymentMethod).Return("", expectedErr).Once()

	resp, err := orderService.PayOrder(ctx, orderUUID, paymentMethod)
	assert.Error(t, err)
	assert.Equal(t, err, expectedErr)
	assert.Empty(t, resp)
}

func TestPayOrderInternalErr(t *testing.T) {
	ctx := context.Background()
	orderUUID := gofakeit.UUID()
	paymentMethod := "CREDIT_CARD"
	expectedErr := model.ErrPaymentInternalError

	order := getMockedPayOrder(orderUUID, model.OrderStatusPendingPayment)

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

	orderRepository.On("GetOrder", mock.Anything, orderUUID).Return(order, nil).Once()
	paymentClient.On("PayOrder", mock.Anything, order.UserUUID, orderUUID, paymentMethod).Return("", expectedErr).Once()

	resp, err := orderService.PayOrder(ctx, orderUUID, paymentMethod)
	assert.Error(t, err)
	assert.Equal(t, err, expectedErr)
	assert.Empty(t, resp)
}

func TestPayOrderNotFoundErr(t *testing.T) {
	ctx := context.Background()
	orderUUID := gofakeit.UUID()
	paymentMethod := "CREDIT_CARD"
	expectedErr := model.ErrPaymentNotFound

	order := getMockedPayOrder(orderUUID, model.OrderStatusPendingPayment)

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

	orderRepository.On("GetOrder", mock.Anything, orderUUID).Return(order, nil).Once()
	paymentClient.On("PayOrder", mock.Anything, order.UserUUID, orderUUID, paymentMethod).Return("", expectedErr).Once()

	resp, err := orderService.PayOrder(ctx, orderUUID, paymentMethod)
	assert.Error(t, err)
	assert.Equal(t, err, expectedErr)
	assert.Empty(t, resp)
}

func TestPayOrderConflictOrderStatusPaidErr(t *testing.T) {
	ctx := context.Background()
	orderUUID := gofakeit.UUID()
	paymentMethod := "CREDIT_CARD"
	expectedErr := model.ErrPaymentConflict

	order := getMockedPayOrder(orderUUID, model.OrderStatusPaid)

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

	orderRepository.On("GetOrder", mock.Anything, orderUUID).Return(order, nil).Once()

	resp, err := orderService.PayOrder(ctx, orderUUID, paymentMethod)
	assert.Error(t, err)
	assert.Equal(t, err, expectedErr)
	assert.Empty(t, resp)
}

func TestPayOrderConflictOrderStatusCanceledErr(t *testing.T) {
	ctx := context.Background()
	orderUUID := gofakeit.UUID()
	paymentMethod := "CREDIT_CARD"
	expectedErr := model.ErrPaymentConflict

	order := getMockedPayOrder(orderUUID, model.OrderStatusCanceled)

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

	orderRepository.On("GetOrder", mock.Anything, orderUUID).Return(order, nil).Once()

	resp, err := orderService.PayOrder(ctx, orderUUID, paymentMethod)
	assert.Error(t, err)
	assert.Equal(t, err, expectedErr)
	assert.Empty(t, resp)
}

func TestPayOrderConflictOrderStatusUnknownErr(t *testing.T) {
	ctx := context.Background()
	orderUUID := gofakeit.UUID()
	paymentMethod := "CREDIT_CARD"
	expectedErr := model.ErrPaymentInternalError

	order := getMockedPayOrder(orderUUID, "UNKNOWN")

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

	orderRepository.On("GetOrder", mock.Anything, orderUUID).Return(order, nil).Once()

	resp, err := orderService.PayOrder(ctx, orderUUID, paymentMethod)
	assert.Error(t, err)
	assert.Equal(t, err, expectedErr)
	assert.Empty(t, resp)
}

func TestPayOrderUpdateErr(t *testing.T) {
	ctx := context.Background()
	orderUUID := gofakeit.UUID()
	paymentMethod := "CREDIT_CARD"
	transactionUUID := gofakeit.UUID()
	expectedErr := model.ErrOrderNotFound

	order := getMockedPayOrder(orderUUID, model.OrderStatusPendingPayment)

	orderInfo := model.OrderUpdateInfo{
		Status:          lo.ToPtr(model.OrderStatusPaid),
		PaymentMethod:   lo.ToPtr(model.PaymentMethod(paymentMethod)),
		TransactionUUID: lo.ToPtr(transactionUUID),
	}

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

	orderRepository.On("GetOrder", mock.Anything, orderUUID).Return(order, nil).Once()
	paymentClient.On("PayOrder", mock.Anything, order.UserUUID, orderUUID, paymentMethod).Return(transactionUUID, nil).Once()
	orderRepository.On("UpdateOrder", mock.Anything, orderUUID, orderInfo).Return(expectedErr).Once()

	_, err := orderService.PayOrder(ctx, orderUUID, paymentMethod)
	assert.Error(t, err)
	assert.Equal(t, err, expectedErr)
}

func getMockedPayOrder(orderUUID string, status model.OrderStatus) model.OrderData {
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
