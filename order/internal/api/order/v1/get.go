package v1

import (
	"context"
	"errors"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/Alexey-step/rocket-factory/order/internal/converter"
	"github.com/Alexey-step/rocket-factory/order/internal/model"
	"github.com/Alexey-step/rocket-factory/platform/pkg/logger"
	orderV1 "github.com/Alexey-step/rocket-factory/shared/pkg/openapi/order/v1"
)

func (a *api) GetOrder(ctx context.Context, params orderV1.GetOrderParams) (orderV1.GetOrderRes, error) {
	order, err := a.service.GetOrder(ctx, params.OrderUUID.String())
	if err != nil {
		if errors.Is(err, model.ErrOrderNotFound) {
			logger.Error(ctx, "Order not found",
				zap.String("order_uuid", params.OrderUUID.String()),
				zap.Error(err),
			)
			return nil, status.Errorf(codes.NotFound, "order by this UUID %s not found", params.OrderUUID.String())
		}
		logger.Error(ctx, "Internal server error while getting order",
			zap.String("order_uuid", params.OrderUUID.String()),
			zap.Error(err),
		)
		return nil, status.Errorf(codes.Internal, "Order service internal error")
	}

	return &orderV1.GetOrderResponse{
		Data: converter.OrderDataToDTO(order),
	}, nil
}
