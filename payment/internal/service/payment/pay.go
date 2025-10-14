package payment

import (
	"context"
	"log"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"github.com/Alexey-step/rocket-factory/platform/pkg/tracing"
)

func (s *service) PayOrder(ctx context.Context, orderUUID, userUUID, paymentMethod string) (transactionUUID string, err error) {
	// Создаем спан для вызова Payment сервиса
	ctx, span := tracing.StartSpan(ctx, "payment.pay_order", //nolint:ineffassign, staticcheck
		trace.WithAttributes(
			attribute.String("order.uuid", orderUUID),
			attribute.String("order.payment_method", paymentMethod),
		),
	)

	defer span.End()

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
