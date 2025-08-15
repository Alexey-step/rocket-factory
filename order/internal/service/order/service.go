package order

import (
	"github.com/Alexey-step/rocket-factory/order/internal/client/grpc"
	def "github.com/Alexey-step/rocket-factory/order/internal/service"
)

var _ def.OrderService = (*service)(nil)

type service struct {
	orderRepository def.OrderRepository

	inventoryClient grpc.InventoryClient
	paymentClient   grpc.PaymentClient

	orderProducerService def.OrderProducerService
}

func NewService(
	orderRepository def.OrderRepository,
	inventoryClient grpc.InventoryClient,
	paymentClient grpc.PaymentClient,
	orderProducerService def.OrderProducerService,
) *service {
	return &service{
		orderRepository:      orderRepository,
		inventoryClient:      inventoryClient,
		paymentClient:        paymentClient,
		orderProducerService: orderProducerService,
	}
}
