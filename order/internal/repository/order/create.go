package order

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/Alexey-step/rocket-factory/order/internal/model"
	repoModel "github.com/Alexey-step/rocket-factory/order/internal/repository/model"
)

func (r *repository) CreateOrder(ctx context.Context, userUUID string, parts []model.Part) (info model.OrderCreationInfo, err error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	orUUID := uuid.New().String()

	var partUUIDs []string
	var totPrice float64
	for _, part := range parts {
		partUUIDs = append(partUUIDs, part.UUID)
		totPrice += part.Price
	}

	order := repoModel.OrderData{
		UUID:       orUUID,
		UserUUID:   userUUID,
		PartUuids:  partUUIDs,
		TotalPrice: totPrice,
		Status:     repoModel.OrderStatusPendingPayment,
		CreatedAt:  time.Now(),
	}

	r.orders[orUUID] = order

	return model.OrderCreationInfo{
		OrderUUID:  orUUID,
		TotalPrice: totPrice,
	}, nil
}
