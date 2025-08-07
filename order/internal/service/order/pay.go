package order

import (
	"context"

	"github.com/samber/lo"
	"go.uber.org/zap"

	"github.com/Alexey-step/rocket-factory/order/internal/model"
	"github.com/Alexey-step/rocket-factory/platform/pkg/logger"
)

func (s *service) PayOrder(ctx context.Context, orderUUID, paymentMethod string) (transactionUUID string, err error) {
	order, err := s.orderRepository.GetOrder(ctx, orderUUID)
	if err != nil {
		logger.Error(ctx, "failed to get order when paying",
			zap.String("order_uuid", orderUUID),
			zap.String("payment_method", paymentMethod),
			zap.Error(err))
		return "", err
	}

	if resp, ok := canPayOrder(order); ok {
		logger.Error(ctx, "failed to pay order",
			zap.String("order_uuid", orderUUID),
			zap.String("payment_method", paymentMethod),
			zap.Error(resp))
		return "", resp
	}

	transUUID, err := s.paymentClient.PayOrder(ctx, order.UserUUID, orderUUID, paymentMethod)
	if err != nil {
		logger.Error(ctx, "failed to pay order",
			zap.String("order_uuid", orderUUID),
			zap.String("payment_method", paymentMethod),
			zap.Error(err))
		return "", err
	}

	orderStatus := model.OrderStatusPaid
	updateErr := s.orderRepository.UpdateOrder(ctx, order.UUID, model.OrderUpdateInfo{
		Status:          &orderStatus,
		PaymentMethod:   lo.ToPtr(model.PaymentMethod(paymentMethod)),
		TransactionUUID: lo.ToPtr(transUUID),
	})

	if updateErr != nil {
		logger.Error(ctx, "failed to update order after payment",
			zap.String("order_uuid", orderUUID),
			zap.String("payment_method", paymentMethod),
			zap.Error(updateErr))
		return "", updateErr
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
