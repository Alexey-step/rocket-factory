package order

import (
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/samber/lo"

	"github.com/Alexey-step/rocket-factory/order/internal/model"
)

func (r *repository) UpdateOrder(ctx context.Context, orderUUID string, orderUpdateInfo model.OrderUpdateInfo) error {
	r.mu.RLock()
	defer r.mu.RUnlock()

	updateBuilder := sq.Update("orders").
		PlaceholderFormat(sq.Dollar).
		Set("updated_at", time.Now())

	// Обновляем поля, только если они были установлены в запросе
	if orderUpdateInfo.Status != nil {
		updateBuilder = updateBuilder.Set("status", lo.FromPtr(orderUpdateInfo.Status))
	}

	if orderUpdateInfo.PaymentMethod != nil {
		updateBuilder = updateBuilder.Set("payment_method", lo.ToPtr(lo.FromPtr(orderUpdateInfo.PaymentMethod)))
	}

	if orderUpdateInfo.TotalPrice != nil {
		updateBuilder = updateBuilder.Set("total_price", *orderUpdateInfo.TotalPrice)
	}

	if orderUpdateInfo.TransactionUUID != nil {
		updateBuilder = updateBuilder.Set("transaction_uuid", orderUpdateInfo.TransactionUUID)
	}

	updateBuilder.Where(sq.Eq{"uuid": orderUUID})

	query, args, err := updateBuilder.ToSql()
	if err != nil {
		return err
	}

	// Выполняем запрос через пул соединений
	_, execErr := r.db.Exec(ctx, query, args...)
	if execErr != nil {
		return execErr
	}

	return nil
}
