package order

import (
	"context"
	"time"

	"github.com/samber/lo"

	"github.com/Alexey-step/rocket-factory/order/internal/model"
	repoModel "github.com/Alexey-step/rocket-factory/order/internal/repository/model"
)

func (r *repository) UpdateOrder(_ context.Context, orderUUID string, orderUpdateInfo model.OrderUpdateInfo) error {
	r.mu.RLock()
	defer r.mu.RUnlock()

	order, ok := r.orders[orderUUID]
	if !ok {
		return model.ErrOrderNotFound
	}

	// Обновляем поля, только если они были установлены в запросе
	if orderUpdateInfo.Status != nil {
		order.Status = repoModel.OrderStatus(lo.FromPtr(orderUpdateInfo.Status))
	}

	if orderUpdateInfo.PaymentMethod != nil {
		order.PaymentMethod = lo.ToPtr(lo.FromPtr(order.PaymentMethod))
	}

	if orderUpdateInfo.TotalPrice != nil {
		order.TotalPrice = *orderUpdateInfo.TotalPrice
	}

	if orderUpdateInfo.TransactionUUID != nil {
		order.TransactionUUID = orderUpdateInfo.TransactionUUID
	}

	order.UpdatedAt = lo.ToPtr(time.Now())

	r.orders[order.UUID] = order

	return nil
}
