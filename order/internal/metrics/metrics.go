package metrics

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
)

var (
	// OrdersTotal - COUNTER для подсчета общего количества заказов
	// Тип: Int64Counter (монотонно возрастающий)
	// Лейблы: method (название метода), status (success/error)
	OrdersTotal metric.Int64Counter

	// OrderRevenueTotal - COUNTER для подсчета общей выручки от заказов
	// Тип: Int64Counter (монотонно возрастающий)
	// Лейблы: method (название метода), status (success/error)
	OrdersRevenueTotal metric.Float64Counter

	// RequestDuration - HISTOGRAM для измерения времени выполнения запросов
	// Тип: Float64Histogram (распределение значений)
	// Использование: SLA мониторинг - отслеживание времени ответа API
	// Позволяет строить percentile (p50, p95, p99) для анализа производительности
	RequestDuration metric.Float64Histogram

	// RequestsTotal - COUNTER для подсчета общего количества запросов
	// Тип: Int64Counter (монотонно возрастающий)
	// Использование: подсчет всех gRPC запросов с разбивкой по методам и статусам
	// Лейблы: method (название метода), status (success/error)
	RequestsTotal metric.Int64Counter
)

type Config interface {
	CollectorServiceName() string
}

// InitMetrics инициализирует все метрики Assembly сервиса
// Должна быть вызвана один раз при старте приложения после инициализации OpenTelemetry провайдера
func InitMetrics(cfg Config) error {
	var err error
	meter := otel.Meter(cfg.CollectorServiceName())

	// Создаем счетчик заказов с описанием для документации
	OrdersTotal, err = meter.Int64Counter(
		"orders_total",
		metric.WithDescription("Total number of orders"),
	)
	if err != nil {
		return err
	}

	// Создаем счетчик выручки с описанием для документации
	OrdersRevenueTotal, err = meter.Float64Counter(
		"orders_revenue_total",
		metric.WithDescription("Total revenue from orders"),
	)
	if err != nil {
		return err
	}

	// Создаем гистограмму времени запросов с правильными bucket'ами для gRPC
	// Bucket'ы оптимизированы для времени отклика в диапазоне от микросекунд до секунд
	RequestDuration, err = meter.Float64Histogram(
		"order_request_duration_seconds",
		metric.WithDescription("Duration of http requests"),
		metric.WithUnit("s"),
		// Добавляем explicit bucket boundaries для более точного измерения gRPC запросов
		// 1ms, 2ms, 5ms, 10ms, 25ms, 50ms, 100ms, 250ms, 500ms, 1s, 2s, 5s
		metric.WithExplicitBucketBoundaries(
			0.001, 0.002, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1.0, 2.0, 5.0,
		),
	)
	if err != nil {
		return err
	}

	// Создаем счетчик запросов с описанием для документации
	RequestsTotal, err = meter.Int64Counter(
		"order_requests_total",
		metric.WithDescription("Total number of Order service requests"),
	)
	if err != nil {
		return err
	}

	return nil
}
