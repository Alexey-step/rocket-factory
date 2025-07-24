package v1

import (
	"context"
	"errors"
	"net/http"

	"github.com/Alexey-step/rocket-factory/order/internal/converter"
	"github.com/Alexey-step/rocket-factory/order/internal/model"
	orderV1 "github.com/Alexey-step/rocket-factory/shared/pkg/openapi/order/v1"
)

func (a *api) PayOrder(ctx context.Context, req *orderV1.PayOrderRequest, params orderV1.PayOrderParams) (orderV1.PayOrderRes, error) {
	transUUID, err := a.service.PayOrder(ctx, params.OrderUUID.String(), string(req.GetPaymentMethod()))
	if err != nil {
		switch {
		case errors.Is(err, model.ErrPaymentConflict):
			return &orderV1.ConflictError{
				Code:    http.StatusBadRequest,
				Message: "Заказ уже оплачен или отменен",
			}, nil
		case errors.Is(err, model.ErrPaymentNotFound):
			return &orderV1.NotFoundError{
				Code:    http.StatusBadRequest,
				Message: "Заказ не найден или не существует",
			}, nil
		default:
			return &orderV1.InternalServerError{
				Code:    http.StatusInternalServerError,
				Message: "Ошибка сервера при проведении платежа",
			}, nil
		}
	}

	return &orderV1.PayOrderResponse{
		TransactionUUID: converter.StringToUUID(transUUID),
	}, nil
}
