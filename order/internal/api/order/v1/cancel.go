package v1

import (
	"context"
	"errors"

	"github.com/Alexey-step/rocket-factory/order/internal/model"
	orderV1 "github.com/Alexey-step/rocket-factory/shared/pkg/openapi/order/v1"
)

func (a *api) CancelOrder(ctx context.Context, params orderV1.CancelOrderParams) (orderV1.CancelOrderRes, error) {
	err := a.service.CancelOrder(ctx, params.OrderUUID.String())
	if err != nil {
		switch {
		case errors.Is(err, model.ErrOrderNotFound):
			return &orderV1.NotFoundError{
				Code:    404,
				Message: "Order by this UUID `" + params.OrderUUID.String() + "` not found",
			}, nil
		case errors.Is(err, model.ErrOrderAlreadyPaid):
			return &orderV1.ConflictError{
				Code:    409,
				Message: "Заказ уже оплачен и не может быть отменён",
			}, nil
		case errors.Is(err, model.ErrOrderAlreadyCancelled):
			return &orderV1.ConflictError{
				Code:    409,
				Message: "Заказ уже отменён",
			}, nil
		default:
			return &orderV1.InternalServerError{
				Code:    500,
				Message: "Внутренняя ошибка сервера",
			}, nil
		}
	}
	return nil, nil
}
