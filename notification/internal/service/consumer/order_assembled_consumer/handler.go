package order_assembled_consumer

import (
	"context"
	"math/rand/v2"
	"time"

	"go.uber.org/zap"

	"github.com/Alexey-step/rocket-factory/notification/internal/model"
	"github.com/Alexey-step/rocket-factory/platform/pkg/kafka"
	"github.com/Alexey-step/rocket-factory/platform/pkg/logger"
)

func (s *service) OrderAssembledHandler(ctx context.Context, msg kafka.Message) error {
	event, err := s.orderAssembledDecoder.Decode(msg.Value)
	if err != nil {
		logger.Error(ctx, "Failed to decode order paid event", zap.Error(err))
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

	//nolint:gosec
	delay := time.Duration(rand.IntN(10)+1) * time.Second
	select {
	case <-time.After(delay):
	case <-ctx.Done():
		return ctx.Err()
	}

	shipAssembled := model.ShipAssembled{
		EventUUID:    event.EventUUID,
		OrderUUID:    event.OrderUUID,
		UserUUID:     event.UserUUID,
		BuildTimeSec: int64(delay / time.Second), // Здесь можно использовать реальное время сборки
	}

	err = s.telegramService.SendAssembledNotification(ctx, shipAssembled)
	if err != nil {
		logger.Error(ctx, "Failed to produce ship assembled event to Telegram",
			zap.Any("ship_assembled", shipAssembled),
			zap.Error(err),
		)
		return err
	}

	return nil
}
