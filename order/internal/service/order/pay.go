package order

import (
	"context"

	"github.com/google/uuid"
	"github.com/samber/lo"

	"github.com/Alexey-step/rocket-factory/order/internal/model"
)

func (s *service) PayOrder(ctx context.Context, orderUUID, paymentMethod string) (transactionUUID string, err error) {
	order, err := s.orderRepository.GetOrder(ctx, orderUUID)
	if err != nil {
		return "", err
	}

	if resp, ok := canPayOrder(order); ok {
		return "", resp
	}

	transUUID, err := s.paymentClient.PayOrder(ctx, order.UserUUID, orderUUID, paymentMethod)
	if err != nil {
		return "", err
	}

	orderStatus := model.OrderStatusPaid
	updateErr := s.orderRepository.UpdateOrder(ctx, order.UUID, model.OrderUpdateInfo{
		Status:          &orderStatus,
		PaymentMethod:   lo.ToPtr(model.PaymentMethod(paymentMethod)),
		TransactionUUID: lo.ToPtr(transUUID),
	})

	if updateErr != nil {
		return "", updateErr
	}

	err = s.orderProducerService.ProduceOrderPaid(ctx, model.OrderPaid{
		EventUUID:       uuid.NewString(),
		OrderUUID:       orderUUID,
		UserUUID:        order.UserUUID,
		PaymentMethod:   paymentMethod,
		TransactionUUID: transUUID,
	})
	if err != nil {
		return "", err
	}

	return transUUID, nil
}

func canPayOrder(order model.OrderData) (error, bool) {
	switch order.Status {
	case model.OrderStatusPaid:
		return model.ErrPaymentConflict, true
	case model.OrderStatusCanceled:
		return model.ErrPaymentConflict, true
	case model.OrderStatusPendingPayment:
		return nil, false
	default:
		return model.ErrPaymentInternalError, true
	}
}
