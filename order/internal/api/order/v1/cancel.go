package v1

import (
	"context"
	"errors"

	"go.uber.org/zap"

	"github.com/Alexey-step/rocket-factory/order/internal/model"
	"github.com/Alexey-step/rocket-factory/platform/pkg/logger"
	orderV1 "github.com/Alexey-step/rocket-factory/shared/pkg/openapi/order/v1"
)

func (a *api) CancelOrder(ctx context.Context, params orderV1.CancelOrderParams) (orderV1.CancelOrderRes, error) {
	err := a.service.CancelOrder(ctx, params.OrderUUID.String())
	if err != nil {
		switch {
		case errors.Is(err, model.ErrOrderNotFound):
			logger.Error(ctx, "Order not found",
				zap.String("order_uuid", params.OrderUUID.String()),
				zap.Error(err),
			)
			return &orderV1.NotFoundError{
				Code:    404,
				Message: "Order by this UUID `" + params.OrderUUID.String() + "` not found",
			}, nil
		case errors.Is(err, model.ErrOrderAlreadyPaid):
			logger.Error(ctx, "Order already paid, cannot cancel",
				zap.String("order_uuid", params.OrderUUID.String()),
				zap.Error(err),
			)
			return &orderV1.ConflictError{
				Code:    409,
				Message: "Заказ уже оплачен и не может быть отменён",
			}, nil
		case errors.Is(err, model.ErrOrderAlreadyCancelled):
			logger.Error(ctx, "Order already cancelled",
				zap.String("order_uuid", params.OrderUUID.String()),
				zap.Error(err),
			)
			return &orderV1.ConflictError{
				Code:    409,
				Message: "Заказ уже отменён",
			}, nil
		default:
			logger.Error(ctx, "Internal server error while cancelling order",
				zap.String("order_uuid", params.OrderUUID.String()),
				zap.Error(err),
			)
			return &orderV1.InternalServerError{
				Code:    500,
				Message: "Внутренняя ошибка сервера",
			}, nil
		}
	}
	return nil, nil
}
