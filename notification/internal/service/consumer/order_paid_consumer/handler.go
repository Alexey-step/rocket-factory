package order_assembled_consumer

import (
	"context"

	"go.uber.org/zap"

	"github.com/Alexey-step/rocket-factory/notification/internal/model"
	"github.com/Alexey-step/rocket-factory/platform/pkg/kafka"
	"github.com/Alexey-step/rocket-factory/platform/pkg/logger"
)

func (s *service) OrderPaidHandler(ctx context.Context, msg kafka.Message) error {
	event, err := s.orderPaidDecoder.Decode(msg.Value)
	if err != nil {
		logger.Error(ctx, "Failed to decode order paid event", zap.Error(err))
		return err
	}

	logger.Info(ctx, "Processing message",
		zap.String("topic", msg.Topic),
		zap.Int32("partition", msg.Partition),
		zap.Int64("offset", msg.Offset),
		zap.String("order_uuid", event.OrderUUID),
		zap.String("event_uuid", event.EventUUID),
		zap.String("payment_method", event.PaymentMethod),
		zap.String("transaction_uuid", event.TransactionUUID),
	)

	orderPaid := model.OrderPaid{
		EventUUID:       event.EventUUID,
		OrderUUID:       event.OrderUUID,
		UserUUID:        event.UserUUID,
		PaymentMethod:   event.PaymentMethod,
		TransactionUUID: event.TransactionUUID,
	}

	err = s.telegramService.SendPaidNotification(ctx, orderPaid)
	if err != nil {
		logger.Error(ctx, "Failed to produce order paid event to Telegram",
			zap.String("event_uuid", orderPaid.EventUUID),
			zap.String("order_uuid", orderPaid.OrderUUID),
			zap.String("payment_method", orderPaid.PaymentMethod),
			zap.String("transaction_uuid", orderPaid.TransactionUUID),
			zap.Error(err),
		)
		return err
	}

	return nil
}
