package payment

import (
	"context"
	"log"

	"github.com/google/uuid"
)

func (s *service) PayOrder(ctx context.Context, orderUUID, userUUID, paymentMethod string) (transactionUUID string, err error) {
	log.Printf(`
💳 [Order Paid]
• 🆔 Order UUID: %s
• 👤 User UUID: %s
• 💰 Payment Method: %s
`, orderUUID, userUUID, paymentMethod,
	)

	res := uuid.New().String()

	log.Printf("✅Оплата прошла успешно, transaction_uuid: %v\n", transactionUUID)

	return res, nil
}
