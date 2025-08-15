package order_consumer

import (
	"context"

	"go.uber.org/zap"

	kafkaConverter "github.com/Alexey-step/rocket-factory/order/internal/converter/kafka"
	serviceOrder "github.com/Alexey-step/rocket-factory/order/internal/service"
	"github.com/Alexey-step/rocket-factory/platform/pkg/kafka"
	"github.com/Alexey-step/rocket-factory/platform/pkg/logger"
)

type service struct {
	orderShipAssembledDecoder kafkaConverter.OrderAssembledDecoder
	orderConsumer             kafka.Consumer
	orderService              serviceOrder.OrderService
}

func NewService(
	orderShipAssembledDecoder kafkaConverter.OrderAssembledDecoder,
	orderConsumer kafka.Consumer,
	orderService serviceOrder.OrderService,
) *service {
	return &service{
		orderShipAssembledDecoder: orderShipAssembledDecoder,
		orderConsumer:             orderConsumer,
		orderService:              orderService,
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
