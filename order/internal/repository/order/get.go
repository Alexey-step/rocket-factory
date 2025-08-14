package order

import (
	"context"

	sq "github.com/Masterminds/squirrel"

	"github.com/Alexey-step/rocket-factory/order/internal/model"
	"github.com/Alexey-step/rocket-factory/order/internal/repository/converter"
	repoModel "github.com/Alexey-step/rocket-factory/order/internal/repository/model"
)

func (r *repository) GetOrder(ctx context.Context, orderUUID string) (model.OrderData, error) {
	query, args, err := sq.Select(
		"uuid",
		"user_uuid",
		"part_uuids",
		"total_price",
		"transaction_uuid",
		"payment_method",
		"status",
		"created_at",
		"updated_at").
		From("orders").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"uuid": orderUUID}).
		ToSql()
	if err != nil {
		return model.OrderData{}, err
	}

	var outOrder repoModel.OrderData
	err = r.db.QueryRow(ctx, query, args...).Scan(
		&outOrder.UUID,
		&outOrder.UserUUID,
		&outOrder.PartUuids,
		&outOrder.TotalPrice,
		&outOrder.TransactionUUID,
		&outOrder.PaymentMethod,
		&outOrder.Status,
		&outOrder.CreatedAt,
		&outOrder.UpdatedAt,
	)
	if err != nil {
		return model.OrderData{}, err
	}

	return converter.OrderDataToModel(outOrder), nil
}
