package order

import (
	"context"

	"github.com/Alexey-step/rocket-factory/order/internal/model"
)

func (s *service) GetOrder(ctx context.Context, orderUUID string) (order model.OrderData, err error) {
	outOrder, err := s.orderRepository.GetOrder(ctx, orderUUID)
	if err != nil {
		return model.OrderData{}, err
	}

	return outOrder, nil
}
