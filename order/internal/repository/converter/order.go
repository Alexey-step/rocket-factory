package converter

import (
	"github.com/samber/lo"

	"github.com/Alexey-step/rocket-factory/order/internal/model"
	repoModel "github.com/Alexey-step/rocket-factory/order/internal/repository/model"
)

func OrderDataToModel(order repoModel.OrderData) model.OrderData {
	return model.OrderData{
		UUID:            order.UUID,
		UserUUID:        order.UserUUID,
		PartUuids:       order.PartUuids,
		TotalPrice:      order.TotalPrice,
		TransactionUUID: order.TransactionUUID,
		PaymentMethod:   lo.ToPtr(model.PaymentMethod(lo.FromPtr(order.PaymentMethod))),
		Status:          model.OrderStatus(order.Status),
		CreatedAt:       order.CreatedAt,
		UpdatedAt:       order.UpdatedAt,
	}
}

func OrderCreateInfoToModel(orderCreateInfo repoModel.OrderCreationInfo) model.OrderCreationInfo {
	return model.OrderCreationInfo{
		OrderUUID:  orderCreateInfo.OrderUUID,
		TotalPrice: orderCreateInfo.TotalPrice,
	}
}
