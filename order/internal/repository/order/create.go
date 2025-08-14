package order

import (
	"context"
	"log"
	"time"

	sq "github.com/Masterminds/squirrel"

	"github.com/Alexey-step/rocket-factory/order/internal/model"
	"github.com/Alexey-step/rocket-factory/order/internal/repository/converter"
	repoModel "github.com/Alexey-step/rocket-factory/order/internal/repository/model"
)

func (r *repository) CreateOrder(ctx context.Context, userUUID string, parts []model.Part) (info model.OrderCreationInfo, err error) {
	var partUUIDs []string
	var totPrice float64
	for _, part := range parts {
		partUUIDs = append(partUUIDs, part.UUID)
		totPrice += part.Price
	}

	order := repoModel.OrderData{
		UserUUID:   userUUID,
		PartUuids:  partUUIDs,
		TotalPrice: totPrice,
		Status:     repoModel.OrderStatusPendingPayment,
		CreatedAt:  time.Now(),
	}

	query, args, err := sq.Insert("orders").
		PlaceholderFormat(sq.Dollar).
		Columns("user_uuid", "part_uuids", "total_price", "status", "created_at").
		Values(order.UserUUID, order.PartUuids, order.TotalPrice, order.Status, order.CreatedAt).
		Suffix("RETURNING uuid, total_price").
		ToSql()
	if err != nil {
		return model.OrderCreationInfo{}, err
	}

	var creationInfo repoModel.OrderCreationInfo
	err = r.db.QueryRow(ctx, query, args...).Scan(&creationInfo.OrderUUID, &creationInfo.TotalPrice)
	if err != nil {
		return model.OrderCreationInfo{}, err
	}

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

	return converter.OrderCreateInfoToModel(creationInfo), nil
}
