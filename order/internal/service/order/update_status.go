package order

import (
	"context"

	"go.uber.org/zap"

	"github.com/Alexey-step/rocket-factory/order/internal/model"
	"github.com/Alexey-step/rocket-factory/platform/pkg/logger"
)

func (s *service) UpdateStatus(ctx context.Context, orderUUID string, status model.OrderStatus) error {
	err := s.orderRepository.UpdateOrder(ctx, orderUUID, model.OrderUpdateInfo{
		Status: &status,
	})
	if err != nil {
		logger.Error(ctx, "update status error",
			zap.String("orderUUID", orderUUID),
			zap.String("status", string(status)),
			zap.Error(err),
		)
		return err
	}

	return nil
}
