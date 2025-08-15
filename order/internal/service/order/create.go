package order

import (
	"context"

	"github.com/Alexey-step/rocket-factory/order/internal/model"
)

func (s *service) CreateOrder(ctx context.Context, userUUID string, partsUUIDs []string) (info model.OrderCreationInfo, error error) {
	filter := model.PartsFilter{
		Uuids: partsUUIDs,
	}

	partsList, err := s.inventoryClient.ListParts(ctx, filter)
	if err != nil {
		return model.OrderCreationInfo{}, err
	}

	if len(partsList) != len(partsUUIDs) {
		return model.OrderCreationInfo{}, model.ErrOrderConflict
	}

	orderInfo, createOrderErr := s.orderRepository.CreateOrder(ctx, userUUID, partsList)
	if createOrderErr != nil {
		return model.OrderCreationInfo{}, createOrderErr
	}

	return model.OrderCreationInfo{
		OrderUUID:  orderInfo.OrderUUID,
		TotalPrice: orderInfo.TotalPrice,
	}, nil
}
