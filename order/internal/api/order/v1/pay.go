package v1

import (
	"context"
	"errors"
	"net/http"

	"go.uber.org/zap"

	"github.com/Alexey-step/rocket-factory/order/internal/converter"
	"github.com/Alexey-step/rocket-factory/order/internal/model"
	"github.com/Alexey-step/rocket-factory/platform/pkg/logger"
	orderV1 "github.com/Alexey-step/rocket-factory/shared/pkg/openapi/order/v1"
)

func (a *api) PayOrder(ctx context.Context, req *orderV1.PayOrderRequest, params orderV1.PayOrderParams) (orderV1.PayOrderRes, error) {
	transUUID, err := a.service.PayOrder(ctx, params.OrderUUID.String(), string(req.GetPaymentMethod()))
	if err != nil {
		switch {
		case errors.Is(err, model.ErrPaymentConflict):
			logger.Error(ctx, "Order payment conflict",
				zap.String("order_uuid", params.OrderUUID.String()),
				zap.String("payment_method", string(req.GetPaymentMethod())),
				zap.Error(err),
			)
			return &orderV1.ConflictError{
				Code:    http.StatusBadRequest,
				Message: "Заказ уже оплачен или отменен",
			}, nil
		case errors.Is(err, model.ErrPaymentNotFound):
			logger.Error(ctx, "Order payment not found",
				zap.String("order_uuid", params.OrderUUID.String()),
				zap.String("payment_method", string(req.GetPaymentMethod())),
				zap.Error(err),
			)
			return &orderV1.NotFoundError{
				Code:    http.StatusBadRequest,
				Message: "Заказ не найден или не существует",
			}, nil
		default:
			logger.Error(ctx, "Internal server error while processing payment",
				zap.String("order_uuid", params.OrderUUID.String()),
				zap.String("payment_method", string(req.GetPaymentMethod())),
				zap.Error(err),
			)
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
