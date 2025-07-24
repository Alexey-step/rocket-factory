package order

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"

	"github.com/Alexey-step/rocket-factory/order/internal/model"
	repoModel "github.com/Alexey-step/rocket-factory/order/internal/repository/model"
)

func (r *repository) CreateOrder(_ context.Context, userUUID string, parts []model.Part) (info model.OrderCreationInfo, err error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	orUUID := uuid.NewString()

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

	log.Printf(`
💳 [Order Created]
• 🆔 Order UUID: %s
• 👤 User UUID: %s
• 💰 Part UUIDs: %v
• 💰 Total Price: %f
• 💰 Status: %s
• 💰 CreatedAt: %v
`, order.UUID, order.UserUUID, order.PartUuids, order.TotalPrice, order.Status, order.CreatedAt,
	)

	return model.OrderCreationInfo{
		OrderUUID:  orUUID,
		TotalPrice: totPrice,
	}, nil
}
