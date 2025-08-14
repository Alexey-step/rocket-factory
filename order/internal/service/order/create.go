package order

import (
	"context"

	"go.uber.org/zap"

	"github.com/Alexey-step/rocket-factory/order/internal/model"
	"github.com/Alexey-step/rocket-factory/platform/pkg/logger"
)

func (s *service) CreateOrder(ctx context.Context, userUUID string, partsUUIDs []string) (info model.OrderCreationInfo, error error) {
	filter := model.PartsFilter{
		Uuids: partsUUIDs,
	}

	partsList, err := s.inventoryClient.ListParts(ctx, filter)
	if err != nil {
		logger.Error(ctx, "failed to get list of parts for order, when creating order",
			zap.Any("filter", filter),
			zap.Strings("parts_uuids", partsUUIDs),
			zap.String("user_uuid", userUUID),
			zap.Error(err),
		)
		return model.OrderCreationInfo{}, err
	}

	if len(partsList) != len(partsUUIDs) {
		logger.Error(ctx, "not all parts found for order",
			zap.Any("filter", filter),
			zap.Strings("parts_uuids", partsUUIDs),
			zap.String("user_uuid", userUUID),
			zap.Any("parts_list", partsList))
		return model.OrderCreationInfo{}, model.ErrOrderConflict
	}

	orderInfo, createOrderErr := s.orderRepository.CreateOrder(ctx, userUUID, partsList)
	if createOrderErr != nil {
		logger.Error(ctx, "failed to create order",
			zap.Any("filter", filter),
			zap.Strings("parts_uuids", partsUUIDs),
			zap.String("user_uuid", userUUID),
			zap.Any("parts_list", partsList))
		return model.OrderCreationInfo{}, createOrderErr
	}

	return model.OrderCreationInfo{
		OrderUUID:  orderInfo.OrderUUID,
		TotalPrice: orderInfo.TotalPrice,
	}, nil
}
