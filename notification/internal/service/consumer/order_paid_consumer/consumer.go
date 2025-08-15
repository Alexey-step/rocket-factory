package order_assembled_consumer

import (
	"context"

	"go.uber.org/zap"

	kafkaConverter "github.com/Alexey-step/rocket-factory/notification/internal/converter/kafka"
	notificationService "github.com/Alexey-step/rocket-factory/notification/internal/service"
	"github.com/Alexey-step/rocket-factory/platform/pkg/kafka"
	"github.com/Alexey-step/rocket-factory/platform/pkg/logger"
)

var _ notificationService.OrderPaidConsumerService = (*service)(nil)

type service struct {
	orderPaidConsumer kafka.Consumer
	orderPaidDecoder  kafkaConverter.OrderPaidDecoder
	telegramService   notificationService.TelegramService
}

func NewService(
	orderPaidConsumer kafka.Consumer,
	orderPaidDecoder kafkaConverter.OrderPaidDecoder,
	telegramService notificationService.TelegramService,
) *service {
	return &service{
		orderPaidConsumer: orderPaidConsumer,
		orderPaidDecoder:  orderPaidDecoder,
		telegramService:   telegramService,
	}
}

func (s *service) RunConsumer(ctx context.Context) error {
	logger.Info(ctx, "Starting order paid consumer service")

	err := s.orderPaidConsumer.Consume(ctx, s.OrderPaidHandler)
	if err != nil {
		logger.Error(ctx, "Failed to start order paid consumer service",
			zap.Error(err),
		)
		return err
	}

	logger.Info(ctx, "OrderPaid consumer service started successfully")
	return nil
}
