package v1

import (
	"context"
	"errors"
	"net/http"

	"github.com/Alexey-step/rocket-factory/order/internal/converter"
	"github.com/Alexey-step/rocket-factory/order/internal/model"
	orderV1 "github.com/Alexey-step/rocket-factory/shared/pkg/openapi/order/v1"
)

func (a *api) CreateOrder(ctx context.Context, req *orderV1.CreateOrderRequest) (orderV1.CreateOrderRes, error) {
	orderInfo, err := a.service.CreateOrder(ctx, req.GetUserUUID().String(), converter.UUIDsToStrings(req.GetPartUuids()))
	if err != nil {
		switch {
		case errors.Is(err, model.ErrPartsNotFound):
			return &orderV1.NotFoundError{
				Code:    http.StatusNotFound,
				Message: "Одна или несколько частей не найдены",
			}, nil
		default:
			return &orderV1.InternalServerError{
				Code:    http.StatusInternalServerError,
				Message: "Внутренняя ошибка сервера при создании заказа",
			}, nil
		}
	}

	return &orderV1.CreateOrderResponse{
		OrderUUID:  converter.StringToUUID(orderInfo.OrderUUID),
		TotalPrice: orderInfo.TotalPrice,
	}, nil
}
