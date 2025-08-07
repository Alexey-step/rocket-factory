package order

import (
	"context"

	"go.uber.org/zap"

	"github.com/Alexey-step/rocket-factory/order/internal/model"
	"github.com/Alexey-step/rocket-factory/platform/pkg/logger"
)

func (s *service) GetOrder(ctx context.Context, orderUUID string) (order model.OrderData, err error) {
	outOrder, err := s.orderRepository.GetOrder(ctx, orderUUID)
	if err != nil {
		logger.Error(ctx, "failed to get order",
			zap.String("order_uuid", orderUUID),
			zap.Error(err))
		return model.OrderData{}, err
	}

	return outOrder, nil
}
