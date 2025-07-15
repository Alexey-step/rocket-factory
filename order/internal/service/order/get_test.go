package order

import (
	"github.com/brianvoe/gofakeit/v7"
	"github.com/samber/lo"

	"github.com/Alexey-step/rocket-factory/order/internal/model"
)

func (s *ServiceSuite) TestGetOrderSuccess() {
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

	s.orderRepository.On("GetOrder", s.ctx, orderUUID).Return(order, nil).Once()
	resp, err := s.service.GetOrder(s.ctx, orderUUID)

	s.NoError(err)
	s.Equal(order, resp)
}

func (s *ServiceSuite) TestGetOrderNotFoundErr() {
	orderUUID := gofakeit.UUID()
	expectedErr := model.ErrOrderNotFound

	s.orderRepository.On("GetOrder", s.ctx, orderUUID).Return(model.OrderData{}, expectedErr).Once()
	resp, err := s.service.GetOrder(s.ctx, orderUUID)

	s.Error(err)
	s.Equal(expectedErr, err)
	s.Empty(resp)
}

func (s *ServiceSuite) TestGetOrderInternalErr() {
	orderUUID := gofakeit.UUID()
	expectedErr := model.ErrOrderInternalError

	s.orderRepository.On("GetOrder", s.ctx, orderUUID).Return(model.OrderData{}, expectedErr).Once()
	resp, err := s.service.GetOrder(s.ctx, orderUUID)

	s.Error(err)
	s.Equal(expectedErr, err)
	s.Empty(resp)
}
