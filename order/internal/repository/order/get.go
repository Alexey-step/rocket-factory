package order

import (
	"context"

	"github.com/Alexey-step/rocket-factory/order/internal/model"
	repoConverter "github.com/Alexey-step/rocket-factory/order/internal/repository/converter"
)

func (r *repository) GetOrder(_ context.Context, orderUUID string) (model.OrderData, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	outOrder, ok := r.orders[orderUUID]
	if !ok {
		return model.OrderData{}, model.ErrOrderNotFound
	}

	return repoConverter.OrderDataToModel(outOrder), nil
}
