package v1

import (
	def "github.com/Alexey-step/rocket-factory/order/internal/client/grpc"
	payment_v1 "github.com/Alexey-step/rocket-factory/shared/pkg/proto/payment/v1"
)

var _ def.PaymentClient = (*client)(nil)

type client struct {
	generatedClient payment_v1.PaymentServiceClient
}

func NewClient(
	generatedClient payment_v1.PaymentServiceClient,
) *client {
	return &client{
		generatedClient: generatedClient,
	}
}
