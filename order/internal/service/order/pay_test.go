package order

import (
	"github.com/brianvoe/gofakeit/v7"
	"github.com/samber/lo"

	"github.com/Alexey-step/rocket-factory/order/internal/model"
)

func (s *ServiceSuite) TestPayOrderSuccess() {
	orderUUID := gofakeit.UUID()
	paymentMethod := "CREDIT_CARD"
	transactionUUID := gofakeit.UUID()

	order := model.OrderData{
		UUID:          orderUUID,
		UserUUID:      gofakeit.UUID(),
		PartUuids:     []string{gofakeit.UUID()},
		TotalPrice:    gofakeit.Price(100, 1000),
		PaymentMethod: lo.ToPtr(model.PaymentMethod("CREDIT_CARD")),
		Status:        model.OrderStatusPendingPayment,
		CreatedAt:     gofakeit.Date(),
	}

	orderInfo := model.OrderUpdateInfo{
		Status:          lo.ToPtr(model.OrderStatusPaid),
		PaymentMethod:   lo.ToPtr(model.PaymentMethod(paymentMethod)),
		TransactionUUID: lo.ToPtr(transactionUUID),
	}

	s.orderRepository.On("GetOrder", s.ctx, orderUUID).Return(order, nil).Once()
	s.paymentClient.On("PayOrder", s.ctx, order.UserUUID, orderUUID, paymentMethod).Return(transactionUUID, nil).Once()
	s.orderRepository.On("UpdateOrder", s.ctx, orderUUID, orderInfo).Return(nil).Once()

	resp, err := s.service.PayOrder(s.ctx, orderUUID, paymentMethod)
	s.NoError(err)
	s.Equal(transactionUUID, resp)
}

func (s *ServiceSuite) TestPayOrderFailGetOrder() {
	orderUUID := gofakeit.UUID()
	paymentMethod := "CREDIT_CARD"
	expectedErr := model.ErrOrderNotFound

	s.orderRepository.On("GetOrder", s.ctx, orderUUID).Return(model.OrderData{}, expectedErr).Once()

	resp, err := s.service.PayOrder(s.ctx, orderUUID, paymentMethod)
	s.Error(err)
	s.Equal(err, expectedErr)
	s.Empty(resp)
}

func (s *ServiceSuite) TestPayOrderFail() {
	orderUUID := gofakeit.UUID()
	paymentMethod := "CREDIT_CARD"
	// transactionUUID := gofakeit.UUID()
	expectedErr := model.ErrPaymentConflict

	order := model.OrderData{
		UUID:          orderUUID,
		UserUUID:      gofakeit.UUID(),
		PartUuids:     []string{gofakeit.UUID()},
		TotalPrice:    gofakeit.Price(100, 1000),
		PaymentMethod: lo.ToPtr(model.PaymentMethod("CREDIT_CARD")),
		Status:        model.OrderStatusPendingPayment,
		CreatedAt:     gofakeit.Date(),
	}

	s.orderRepository.On("GetOrder", s.ctx, orderUUID).Return(order, nil).Once()
	s.paymentClient.On("PayOrder", s.ctx, order.UserUUID, orderUUID, paymentMethod).Return("", expectedErr).Once()

	resp, err := s.service.PayOrder(s.ctx, orderUUID, paymentMethod)
	s.Error(err)
	s.Equal(err, expectedErr)
	s.Empty(resp)
}

func (s *ServiceSuite) TestPayOrderInternalErr() {
	orderUUID := gofakeit.UUID()
	paymentMethod := "CREDIT_CARD"
	expectedErr := model.ErrPaymentInternalError

	order := model.OrderData{
		UUID:          orderUUID,
		UserUUID:      gofakeit.UUID(),
		PartUuids:     []string{gofakeit.UUID()},
		TotalPrice:    gofakeit.Price(100, 1000),
		PaymentMethod: lo.ToPtr(model.PaymentMethod("CREDIT_CARD")),
		Status:        model.OrderStatusPendingPayment,
		CreatedAt:     gofakeit.Date(),
	}

	s.orderRepository.On("GetOrder", s.ctx, orderUUID).Return(order, nil).Once()
	s.paymentClient.On("PayOrder", s.ctx, order.UserUUID, orderUUID, paymentMethod).Return("", expectedErr).Once()

	resp, err := s.service.PayOrder(s.ctx, orderUUID, paymentMethod)
	s.Error(err)
	s.Equal(err, expectedErr)
	s.Empty(resp)
}

func (s *ServiceSuite) TestPayOrderNotFoundErr() {
	orderUUID := gofakeit.UUID()
	paymentMethod := "CREDIT_CARD"
	expectedErr := model.ErrPaymentNotFound

	order := model.OrderData{
		UUID:          orderUUID,
		UserUUID:      gofakeit.UUID(),
		PartUuids:     []string{gofakeit.UUID()},
		TotalPrice:    gofakeit.Price(100, 1000),
		PaymentMethod: lo.ToPtr(model.PaymentMethod("CREDIT_CARD")),
		Status:        model.OrderStatusPendingPayment,
		CreatedAt:     gofakeit.Date(),
	}

	s.orderRepository.On("GetOrder", s.ctx, orderUUID).Return(order, nil).Once()
	s.paymentClient.On("PayOrder", s.ctx, order.UserUUID, orderUUID, paymentMethod).Return("", expectedErr).Once()

	resp, err := s.service.PayOrder(s.ctx, orderUUID, paymentMethod)
	s.Error(err)
	s.Equal(err, expectedErr)
	s.Empty(resp)
}

func (s *ServiceSuite) TestPayOrderConflictOrderStatusPaidErr() {
	orderUUID := gofakeit.UUID()
	paymentMethod := "CREDIT_CARD"
	expectedErr := model.ErrPaymentConflict

	order := model.OrderData{
		UUID:          orderUUID,
		UserUUID:      gofakeit.UUID(),
		PartUuids:     []string{gofakeit.UUID()},
		TotalPrice:    gofakeit.Price(100, 1000),
		PaymentMethod: lo.ToPtr(model.PaymentMethod("CREDIT_CARD")),
		Status:        model.OrderStatusPaid,
		CreatedAt:     gofakeit.Date(),
	}

	s.orderRepository.On("GetOrder", s.ctx, orderUUID).Return(order, nil).Once()

	resp, err := s.service.PayOrder(s.ctx, orderUUID, paymentMethod)
	s.Error(err)
	s.Equal(err, expectedErr)
	s.Empty(resp)
}

func (s *ServiceSuite) TestPayOrderConflictOrderStatusCanceledErr() {
	orderUUID := gofakeit.UUID()
	paymentMethod := "CREDIT_CARD"
	expectedErr := model.ErrPaymentConflict

	order := model.OrderData{
		UUID:          orderUUID,
		UserUUID:      gofakeit.UUID(),
		PartUuids:     []string{gofakeit.UUID()},
		TotalPrice:    gofakeit.Price(100, 1000),
		PaymentMethod: lo.ToPtr(model.PaymentMethod("CREDIT_CARD")),
		Status:        model.OrderStatusCanceled,
		CreatedAt:     gofakeit.Date(),
	}

	s.orderRepository.On("GetOrder", s.ctx, orderUUID).Return(order, nil).Once()

	resp, err := s.service.PayOrder(s.ctx, orderUUID, paymentMethod)
	s.Error(err)
	s.Equal(err, expectedErr)
	s.Empty(resp)
}

func (s *ServiceSuite) TestPayOrderConflictOrderStatusUnknownErr() {
	orderUUID := gofakeit.UUID()
	paymentMethod := "CREDIT_CARD"
	expectedErr := model.ErrPaymentInternalError

	order := model.OrderData{
		UUID:          orderUUID,
		UserUUID:      gofakeit.UUID(),
		PartUuids:     []string{gofakeit.UUID()},
		TotalPrice:    gofakeit.Price(100, 1000),
		PaymentMethod: lo.ToPtr(model.PaymentMethod("CREDIT_CARD")),
		Status:        "UNKNOWN",
		CreatedAt:     gofakeit.Date(),
	}

	s.orderRepository.On("GetOrder", s.ctx, orderUUID).Return(order, nil).Once()

	resp, err := s.service.PayOrder(s.ctx, orderUUID, paymentMethod)
	s.Error(err)
	s.Equal(err, expectedErr)
	s.Empty(resp)
}

func (s *ServiceSuite) TestPayOrderUpdateErr() {
	orderUUID := gofakeit.UUID()
	paymentMethod := "CREDIT_CARD"
	transactionUUID := gofakeit.UUID()
	expectedErr := model.ErrOrderNotFound

	order := model.OrderData{
		UUID:          orderUUID,
		UserUUID:      gofakeit.UUID(),
		PartUuids:     []string{gofakeit.UUID()},
		TotalPrice:    gofakeit.Price(100, 1000),
		PaymentMethod: lo.ToPtr(model.PaymentMethod("CREDIT_CARD")),
		Status:        model.OrderStatusPendingPayment,
		CreatedAt:     gofakeit.Date(),
	}

	orderInfo := model.OrderUpdateInfo{
		Status:          lo.ToPtr(model.OrderStatusPaid),
		PaymentMethod:   lo.ToPtr(model.PaymentMethod(paymentMethod)),
		TransactionUUID: lo.ToPtr(transactionUUID),
	}

	s.orderRepository.On("GetOrder", s.ctx, orderUUID).Return(order, nil).Once()
	s.paymentClient.On("PayOrder", s.ctx, order.UserUUID, orderUUID, paymentMethod).Return(transactionUUID, nil).Once()
	s.orderRepository.On("UpdateOrder", s.ctx, orderUUID, orderInfo).Return(expectedErr).Once()

	_, err := s.service.PayOrder(s.ctx, orderUUID, paymentMethod)
	s.Error(err)
	s.Equal(err, expectedErr)
}
