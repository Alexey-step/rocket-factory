package order_assembled_consumer

import (
	"context"

	"go.uber.org/zap"

	kafkaConverter "github.com/Alexey-step/rocket-factory/notification/internal/converter/kafka"
	notificationService "github.com/Alexey-step/rocket-factory/notification/internal/service"
	"github.com/Alexey-step/rocket-factory/platform/pkg/kafka"
	"github.com/Alexey-step/rocket-factory/platform/pkg/logger"
)

var _ notificationService.OrderAssembledConsumerService = (*service)(nil)

type service struct {
	orderAssembledConsumer kafka.Consumer
	orderAssembledDecoder  kafkaConverter.OrderAssembledDecoder
	telegramService        notificationService.TelegramService
}

func NewService(
	orderAssembledConsumer kafka.Consumer,
	orderAssembledDecoder kafkaConverter.OrderAssembledDecoder,
	telegramService notificationService.TelegramService,
) *service {
	return &service{
		orderAssembledConsumer: orderAssembledConsumer,
		orderAssembledDecoder:  orderAssembledDecoder,
		telegramService:        telegramService,
	}
}

func (s *service) RunConsumer(ctx context.Context) error {
	logger.Info(ctx, "Starting order assembled consumer service")

	err := s.orderAssembledConsumer.Consume(ctx, s.OrderAssembledHandler)
	if err != nil {
		logger.Error(ctx, "Failed to start order assembled consumer service",
			zap.Error(err),
		)
		return err
	}

	logger.Info(ctx, "OrderAssembled consumer service started successfully")
	return nil
}
