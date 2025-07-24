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

	switch order.Status {
	case model.OrderStatusPaid:
		return model.ErrOrderConflict
	case model.OrderStatusCanceled:
		return model.ErrOrderConflict
	case model.OrderStatusPendingPayment:
		status := model.OrderStatusCanceled
		err = s.orderRepository.UpdateOrder(ctx, order.UUID, model.OrderUpdateInfo{
			Status: &status,
		})
		if err != nil {
			return err
		}
		return nil
	default:
		return model.ErrOrderInternalError
	}
}
