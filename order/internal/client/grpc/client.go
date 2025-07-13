package grpc

import (
	"context"

	"github.com/Alexey-step/rocket-factory/order/internal/model"
)

type InventoryClient interface {
	ListParts(ctx context.Context, filter model.PartsFilter) (parts []model.Part, err error)
}

type PaymentClient interface {
	PayOrder(ctx context.Context, userUUID, orderUUID, paymentMethod string) (transactionUUID string, err error)
}
