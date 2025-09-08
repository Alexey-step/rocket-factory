package metrics

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
)

// RocketBuildingDuration - HISTOGRAM для измерения времени сборки коробля
// Тип: Float64Histogram (распределение значений)
// Использование: SLA мониторинг - отслеживание времени ответа API
// Позволяет строить percentile (p50, p95, p99) для анализа производительности
var RocketBuildingDuration metric.Float64Histogram

type Config interface {
	CollectorServiceName() string
}

// InitMetrics инициализирует все метрики Assembly сервиса
// Должна быть вызвана один раз при старте приложения после инициализации OpenTelemetry провайдера
func InitMetrics(cfg Config) error {
	var err error
	meter := otel.Meter(cfg.CollectorServiceName())
	// Создаем гистограмму времени запросов с правильными bucket'ами для gRPC
	// Bucket'ы оптимизированы для времени отклика в диапазоне от микросекунд до секунд
	RocketBuildingDuration, err = meter.Float64Histogram(
		"assembly_duration_seconds",
		metric.WithDescription("Duration of rocket building"),
		metric.WithUnit("s"),
		// Добавляем explicit bucket boundaries для более точного измерения gRPC запросов
		// 1ms, 2ms, 5ms, 10ms, 25ms, 50ms, 100ms, 250ms, 500ms, 1s, 2s, 5s
		metric.WithExplicitBucketBoundaries(
			1.0, 2.0, 3.0, 5.0, 8.0, 10.0, 15.0, 30.0, 60.0,
		),
	)
	if err != nil {
		return err
	}

	return nil
}
