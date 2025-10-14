package middlewares

import (
	"log"
	"net/http"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"

	orderMetrics "github.com/Alexey-step/rocket-factory/order/internal/metrics"
)

// MetricsMiddleware создает middleware для логирования времени выполнения запросов
func MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Засекаем время начала обработки запроса
		startTime := time.Now()

		// Передаем управление следующему обработчику
		next.ServeHTTP(w, r)

		// Записываем время выполнения
		duration := time.Since(startTime)
		durationSeconds := duration.Seconds()

		log.Printf("🕐 Request duration: %v (%f seconds) for method: %s", duration, durationSeconds, r.Method)

		orderMetrics.RequestDuration.Record(r.Context(), durationSeconds,
			metric.WithAttributes(
				attribute.String("method", r.Method),
			),
		)

		orderMetrics.RequestsTotal.Add(r.Context(), 1,
			metric.WithAttributes(
				attribute.String("method", r.Method),
			),
		)
	})
}
