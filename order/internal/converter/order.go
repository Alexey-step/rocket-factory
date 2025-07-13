package converter

import (
	"log"

	"github.com/google/uuid"
	"github.com/samber/lo"

	"github.com/Alexey-step/rocket-factory/order/internal/model"
	repoModel "github.com/Alexey-step/rocket-factory/order/internal/repository/model"
	orderV1 "github.com/Alexey-step/rocket-factory/shared/pkg/openapi/order/v1"
)

func OrderDataToRepoModel(order model.OrderData) repoModel.OrderData {
	return repoModel.OrderData{
		UUID:            order.UUID,
		UserUUID:        order.UserUUID,
		PartUuids:       order.PartUuids,
		TotalPrice:      order.TotalPrice,
		TransactionUUID: order.TransactionUUID,
		PaymentMethod:   lo.ToPtr(repoModel.PaymentMethod(lo.FromPtr(order.PaymentMethod))),
		Status:          repoModel.OrderStatus(order.Status),
		CreatedAt:       order.CreatedAt,
		UpdatedAt:       order.UpdatedAt,
	}
}

func OrderDataToDTO(order model.OrderData) orderV1.OrderDto {
	var transactionUUID orderV1.OptUUID
	if order.TransactionUUID != nil {
		transactionUUID = orderV1.OptUUID{Value: StringToUUID(*order.TransactionUUID)}
	}

	var paymentMethod orderV1.OptPaymentMethod
	if order.PaymentMethod != nil {
		paymentMethod = orderV1.OptPaymentMethod{Value: paymentMethodToOpt(*order.PaymentMethod)}
	}

	var updatedAt orderV1.OptDateTime
	if order.UpdatedAt != nil {
		updatedAt = orderV1.OptDateTime{Value: *order.UpdatedAt}
	}

	return orderV1.OrderDto{
		OrderUUID:       StringToUUID(order.UUID),
		UserUUID:        StringToUUID(order.UserUUID),
		PartUuids:       stringsToUUIDs(order.PartUuids),
		TotalPrice:      order.TotalPrice,
		TransactionUUID: transactionUUID,
		PaymentMethod:   paymentMethod,
		Status:          orderV1.OrderStatus(order.Status),
		CreatedAt:       order.CreatedAt,
		UpdatedAt:       updatedAt,
	}
}

func StringToUUID(s string) uuid.UUID {
	u, err := uuid.Parse(s)
	if err != nil {
		log.Printf("Failed to parse UUID: %v", err)
	}

	return u
}

func stringsToUUIDs(arr []string) []uuid.UUID {
	uuids := make([]uuid.UUID, len(arr))
	for i, s := range arr {
		uuids[i] = StringToUUID(s)
	}
	return uuids
}

func UUIDsToStrings(arr []uuid.UUID) []string {
	uuids := make([]string, len(arr))
	for i, s := range arr {
		uuids[i] = s.String()
	}
	return uuids
}

func paymentMethodToOpt(pm model.PaymentMethod) orderV1.PaymentMethod {
	// Реализуйте по аналогии с OptUUID, если нужно
	return orderV1.PaymentMethod(pm)
}
