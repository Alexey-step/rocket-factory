package order

import (
	"context"

	"github.com/Alexey-step/rocket-factory/order/internal/model"
)

func (s *service) CancelOrder(ctx context.Context, userUUID string) error {
	order, err := s.orderRepository.GetOrder(ctx, userUUID)
	if err != nil {
		return err
	}

	if err = statusToError(order.Status); err != nil {
		return err
	}

	status := model.OrderStatusCanceled
	err = s.orderRepository.UpdateOrder(ctx, order.UUID, model.OrderUpdateInfo{
		Status: &status,
	})
	if err != nil {
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
