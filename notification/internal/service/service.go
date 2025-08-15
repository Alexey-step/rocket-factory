package service

import (
	"context"

	"github.com/Alexey-step/rocket-factory/notification/internal/model"
)

type OrderPaidConsumerService interface {
	RunConsumer(ctx context.Context) error
}

type OrderAssembledConsumerService interface {
	RunConsumer(ctx context.Context) error
}

type TelegramService interface {
	SendPaidNotification(ctx context.Context, paid model.OrderPaid) error
	SendAssembledNotification(ctx context.Context, assembled model.ShipAssembled) error
}
