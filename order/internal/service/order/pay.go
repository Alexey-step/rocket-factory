package order

import (
	"context"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	orderMetrics "github.com/Alexey-step/rocket-factory/order/internal/metrics"
	"github.com/Alexey-step/rocket-factory/order/internal/model"
	"github.com/Alexey-step/rocket-factory/platform/pkg/tracing"
)

func (s *service) PayOrder(ctx context.Context, orderUUID, paymentMethod string) (transactionUUID string, err error) {
	// Создаем спан для вызова Payment сервиса
	ctx, span := tracing.StartSpan(ctx, "order.call_payment_pay_order",
		trace.WithAttributes(
			attribute.String("order.uuid", orderUUID),
			attribute.String("order.payment_method", paymentMethod),
		),
	)

	defer span.End()

	order, err := s.orderRepository.GetOrder(ctx, orderUUID)
	if err != nil {
		span.RecordError(err)
		return "", err
	}

	if resp, ok := canPayOrder(order); ok {
		return "", resp
	}

	transUUID, err := s.paymentClient.PayOrder(ctx, order.UserUUID, orderUUID, paymentMethod)
	if err != nil {
		span.RecordError(err)
		return "", err
	}

	// Добавляем атрибуты результата
	span.SetAttributes(
		attribute.String("payment.transaction_uuid", transUUID),
	)

	orderStatus := model.OrderStatusPaid
	updateErr := s.orderRepository.UpdateOrder(ctx, order.UUID, model.OrderUpdateInfo{
		Status:          &orderStatus,
		PaymentMethod:   lo.ToPtr(model.PaymentMethod(paymentMethod)),
		TransactionUUID: lo.ToPtr(transUUID),
	})

	if updateErr != nil {
		span.RecordError(err)
		return "", updateErr
	}

	err = s.orderProducerService.ProduceOrderPaid(ctx, model.OrderPaid{
		EventUUID:       uuid.NewString(),
		OrderUUID:       orderUUID,
		UserUUID:        order.UserUUID,
		PaymentMethod:   paymentMethod,
		TransactionUUID: transUUID,
	})
	if err != nil {
		span.RecordError(err)
		return "", err
	}

	// Метрики по общей сумме оплат
	if orderMetrics.OrdersRevenueTotal != nil {
		orderMetrics.OrdersRevenueTotal.Add(ctx, order.TotalPrice)
	}

	return transUUID, nil
}

func canPayOrder(order model.OrderData) (error, bool) {
	switch order.Status {
	case model.OrderStatusPaid:
		return model.ErrPaymentConflict, true
	case model.OrderStatusCanceled:
		return model.ErrPaymentConflict, true
	case model.OrderStatusPendingPayment:
		return nil, false
	default:
		return model.ErrPaymentInternalError, true
	}
}
