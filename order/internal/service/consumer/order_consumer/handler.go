package order_consumer

import (
	"context"

	"go.uber.org/zap"

	"github.com/Alexey-step/rocket-factory/order/internal/model"
	"github.com/Alexey-step/rocket-factory/platform/pkg/kafka"
	"github.com/Alexey-step/rocket-factory/platform/pkg/logger"
)

func (s *service) OrderHandler(ctx context.Context, msg kafka.Message) error {
	event, err := s.orderShipAssembledDecoder.Decode(msg.Value)
	if err != nil {
		logger.Error(ctx, "Failed to decode ship assembled event", zap.Error(err))
		return err
	}

	logger.Info(ctx, "Processing message",
		zap.String("topic", msg.Topic),
		zap.Any("partition", msg.Partition),
		zap.Any("offset", msg.Offset),
		zap.String("order_uuid", event.OrderUUID),
		zap.String("event_uuid", event.EventUUID),
		zap.Int64("build_time_sec", event.BuildTimeSec),
	)

	err = s.orderService.UpdateStatus(ctx, event.OrderUUID, model.OrderStatusCompleted)
	if err != nil {
		logger.Error(ctx, "Failed to update order status to completed",
			zap.String("order_uuid", event.OrderUUID),
			zap.Error(err),
		)
		return err
	}

	return nil
}
