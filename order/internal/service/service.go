package service

import (
	"context"

	"github.com/Alexey-step/rocket-factory/order/internal/model"
)

type OrderService interface {
	CreateOrder(ctx context.Context, userUUID string, partsUUIDs []string) (info model.OrderCreationInfo, err error)
	GetOrder(ctx context.Context, orderUUID string) (order model.OrderData, err error)
	CancelOrder(ctx context.Context, orderUUID string) error
	PayOrder(ctx context.Context, orderUUID, paymentMethod string) (transactionUUID string, err error)
}

type OrderRepository interface {
	CreateOrder(ctx context.Context, userUUID string, parts []model.Part) (info model.OrderCreationInfo, err error)
	GetOrder(ctx context.Context, orderUUID string) (order model.OrderData, err error)
	UpdateOrder(ctx context.Context, orderUUID string, orderUpdateInfo model.OrderUpdateInfo) error
}
