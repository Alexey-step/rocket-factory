package order_producer

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"

	"github.com/Alexey-step/rocket-factory/assembly/internal/model"
	"github.com/Alexey-step/rocket-factory/platform/pkg/kafka"
	"github.com/Alexey-step/rocket-factory/platform/pkg/logger"
	eventsV1 "github.com/Alexey-step/rocket-factory/shared/pkg/proto/events/v1"
)

type service struct {
	orderProducer kafka.Producer
}

func NewService(orderProducer kafka.Producer) *service {
	return &service{
		orderProducer: orderProducer,
	}
}

func (s *service) ProduceShipAssembled(ctx context.Context, event model.ShipAssembled) error {
	msg := &eventsV1.ShipAssembled{
		EventUuid:    event.EventUUID,
		OrderUuid:    event.OrderUUID,
		UserUuid:     event.UserUUID,
		BuildTimeSec: event.BuildTimeSec,
	}

	payload, err := proto.Marshal(msg)
	if err != nil {
		logger.Error(ctx, "Failed to marshal order paid", zap.Error(err))
		return err
	}

	err = s.orderProducer.Send(ctx, event.OrderUUID, payload)
	if err != nil {
		logger.Error(ctx, "Failed to send order paid",
			zap.Any("event", event),
			zap.Error(err),
		)
		return err
	}

	return nil
}
