package order_consumer

import (
	"context"

	"go.uber.org/zap"

	kafkaConverter "github.com/Alexey-step/rocket-factory/assembly/internal/converter/kafka"
	assemblyService "github.com/Alexey-step/rocket-factory/assembly/internal/service"
	"github.com/Alexey-step/rocket-factory/platform/pkg/kafka"
	"github.com/Alexey-step/rocket-factory/platform/pkg/logger"
)

var _ assemblyService.OrderConsumerService = (*service)(nil)

type service struct {
	orderConsumer    kafka.Consumer
	orderPaidDecoder kafkaConverter.OrderPaidDecoder
	orderProducer    assemblyService.OrderProducerService
}

func NewService(
	orderConsumer kafka.Consumer,
	orderPaidDecoder kafkaConverter.OrderPaidDecoder,
	orderProducer assemblyService.OrderProducerService,
) *service {
	return &service{
		orderConsumer:    orderConsumer,
		orderPaidDecoder: orderPaidDecoder,
		orderProducer:    orderProducer,
	}
}

func (s *service) RunConsumer(ctx context.Context) error {
	logger.Info(ctx, "Starting order consumer service")

	err := s.orderConsumer.Consume(ctx, s.OrderHandler)
	if err != nil {
		logger.Error(ctx, "Failed to start order consumer service",
			zap.Error(err),
		)
		return err
	}

	logger.Info(ctx, "Order consumer service started successfully")
	return nil
}
