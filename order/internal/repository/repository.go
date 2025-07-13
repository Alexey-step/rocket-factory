package repository

import (
	"context"

	"github.com/Alexey-step/rocket-factory/order/internal/model"
)

type OrderRepository interface {
	CreateOrder(ctx context.Context, userUUID string, parts []model.Part) (info model.OrderCreationInfo, err error)
	GetOrder(ctx context.Context, orderUUID string) (order model.OrderData, err error)
	UpdateOrder(ctx context.Context, orderUUID string, orderUpdateInfo model.OrderUpdateInfo) error
}
