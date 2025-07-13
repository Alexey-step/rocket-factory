package v1

import (
	"context"

	generatedPaymentV1 "github.com/Alexey-step/rocket-factory/shared/pkg/proto/payment/v1"
)

func (c *client) PayOrder(ctx context.Context, userUUID, orderUUID, paymentMethod string) (transactionUUID string, err error) {
	res, err := c.generatedClient.PayOrder(ctx, &generatedPaymentV1.PayOrderRequest{
		OrderUuid:     orderUUID,
		UserUuid:      userUUID,
		PaymentMethod: generatedPaymentV1.PaymentMethod(generatedPaymentV1.PaymentMethod_value[paymentMethod]),
	})
	if err != nil {
		return "", err
	}
	return res.TransactionUuid, nil
}
