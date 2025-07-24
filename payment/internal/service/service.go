package service

import "context"

type PaymentService interface {
	PayOrder(ctx context.Context, orderUUID, userUUID, paymentMethod string) (transactionUUID string, err error)
}
