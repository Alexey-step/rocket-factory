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

func (a *api) CreateOrder(ctx context.Context, req *orderV1.CreateOrderRequest) (orderV1.CreateOrderRes, error) {
	orderInfo, err := a.service.CreateOrder(ctx, req.GetUserUUID().String(), converter.UUIDsToStrings(req.GetPartUuids()))
	if err != nil {
		if errors.Is(err, model.ErrPartsNotFound) {
			logger.Error(ctx, "Some parts not found",
				zap.Strings("part_uuids", converter.UUIDsToStrings(req.GetPartUuids())),
				zap.Error(err),
			)
			return &orderV1.NotFoundError{
				Code:    http.StatusNotFound,
				Message: "Одна или несколько частей не найдены",
			}, nil
		}
		logger.Error(ctx, "Failed to create order",
			zap.Strings("part_uuids", converter.UUIDsToStrings(req.GetPartUuids())),
			zap.Error(err),
		)
		return nil, err
	}

	return &orderV1.CreateOrderResponse{
		OrderUUID:  converter.StringToUUID(orderInfo.OrderUUID),
		TotalPrice: orderInfo.TotalPrice,
	}, nil
}
