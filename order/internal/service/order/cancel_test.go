package order

import (
	"github.com/brianvoe/gofakeit/v7"
	"github.com/samber/lo"

	"github.com/Alexey-step/rocket-factory/order/internal/model"
)

func (s *ServiceSuite) TestCancelOrderSuccess() {
	orderUUID := gofakeit.UUID()

	order := model.OrderData{
		UUID:          orderUUID,
		UserUUID:      gofakeit.UUID(),
		PartUuids:     []string{gofakeit.UUID()},
		TotalPrice:    gofakeit.Price(100, 1000),
		PaymentMethod: lo.ToPtr(model.PaymentMethod("CREDIT_CARD")),
		Status:        model.OrderStatusPendingPayment,
		CreatedAt:     gofakeit.Date(),
	}

	orderUpdateInfo := model.OrderUpdateInfo{
		Status: lo.ToPtr(model.OrderStatusCanceled),
	}

	s.orderRepository.On("GetOrder", s.ctx, orderUUID).Return(order, nil).Once()
	s.orderRepository.On("UpdateOrder", s.ctx, orderUUID, orderUpdateInfo).Return(nil).Once()
	err := s.service.CancelOrder(s.ctx, orderUUID)
	s.NoError(err)
}

func (s *ServiceSuite) TestCancelOrderFail() {
	orderUUID := gofakeit.UUID()
	expectedErr := model.ErrOrderConflict

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
	err := s.service.CancelOrder(s.ctx, orderUUID)
	s.Error(err)
	s.Equal(expectedErr, err)
}

func (s *ServiceSuite) TestCancelOrderConflictFail() {
	orderUUID := gofakeit.UUID()
	expectedErr := model.ErrOrderConflict

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
	err := s.service.CancelOrder(s.ctx, orderUUID)
	s.Error(err)
	s.Equal(expectedErr, err)
}

func (s *ServiceSuite) TestCancelOrderInternalError() {
	orderUUID := gofakeit.UUID()
	expectedErr := model.ErrOrderInternalError

	s.orderRepository.On("GetOrder", s.ctx, orderUUID).Return(model.OrderData{}, expectedErr).Once()
	err := s.service.CancelOrder(s.ctx, orderUUID)
	s.Error(err)
	s.Equal(expectedErr, err)
}

func (s *ServiceSuite) TestCancelOrderNotFoundFail() {
	orderUUID := gofakeit.UUID()
	expectedErr := model.ErrOrderNotFound

	s.orderRepository.On("GetOrder", s.ctx, orderUUID).Return(model.OrderData{}, expectedErr).Once()
	err := s.service.CancelOrder(s.ctx, orderUUID)
	s.Error(err)
	s.Equal(expectedErr, err)
}

func (s *ServiceSuite) TestCancelOrderUpdateError() {
	orderUUID := gofakeit.UUID()
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

	orderUpdateInfo := model.OrderUpdateInfo{
		Status: lo.ToPtr(model.OrderStatusCanceled),
	}

	s.orderRepository.On("GetOrder", s.ctx, orderUUID).Return(order, nil).Once()
	s.orderRepository.On("UpdateOrder", s.ctx, orderUUID, orderUpdateInfo).Return(expectedErr).Once()
	err := s.service.CancelOrder(s.ctx, orderUUID)
	s.Error(err)
	s.Equal(expectedErr, err)
}

func (s *ServiceSuite) TestCancelOrderInternalErr() {
	orderUUID := gofakeit.UUID()
	expectedErr := model.ErrOrderInternalError

	order := model.OrderData{
		UUID:          orderUUID,
		UserUUID:      gofakeit.UUID(),
		PartUuids:     []string{gofakeit.UUID()},
		TotalPrice:    gofakeit.Price(100, 1000),
		PaymentMethod: lo.ToPtr(model.PaymentMethod("CREDIT_CARD")),
		Status:        model.OrderStatus("UNKNOWN_STATUS"),
		CreatedAt:     gofakeit.Date(),
	}

	s.orderRepository.On("GetOrder", s.ctx, orderUUID).Return(order, expectedErr).Once()
	err := s.service.CancelOrder(s.ctx, orderUUID)
	s.Error(err)
	s.Equal(expectedErr, err)
}
