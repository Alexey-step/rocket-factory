package order

import (
	"context"

	"go.uber.org/zap"

	"github.com/Alexey-step/rocket-factory/order/internal/model"
	"github.com/Alexey-step/rocket-factory/platform/pkg/logger"
)

func (s *service) CancelOrder(ctx context.Context, userUUID string) error {
	order, err := s.orderRepository.GetOrder(ctx, userUUID)
	if err != nil {
		logger.Error(ctx, "failed to get order when canceling",
			zap.String("order_uuid", userUUID),
			zap.Error(err),
		)
		return err
	}

	if err = statusToError(order.Status); err != nil {
		logger.Error(ctx, "failed to canceling order",
			zap.String("order_uuid", userUUID),
			zap.Error(err),
		)
		return err
	}

	status := model.OrderStatusCanceled
	err = s.orderRepository.UpdateOrder(ctx, order.UUID, model.OrderUpdateInfo{
		Status: &status,
	})
	if err != nil {
		logger.Error(ctx, "failed to update order when canceling",
			zap.String("order_uuid", userUUID),
			zap.Error(err),
		)
		return err
	}

	return nil
}

func statusToError(status model.OrderStatus) error {
	statusToError := map[model.OrderStatus]error{
		model.OrderStatusPaid:     model.ErrOrderConflict,
		model.OrderStatusCanceled: model.ErrOrderConflict,
	}

	if err, ok := statusToError[status]; ok {
		return err
	}

	if status != model.OrderStatusPendingPayment {
		return model.ErrOrderInternalError
	}

	return nil
}
