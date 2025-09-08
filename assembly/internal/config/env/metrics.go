package env

import (
	"time"

	"github.com/caarlos0/env/v11"
)

type metricsEnvConfig struct {
	OtelServiceName        string        `env:"METRICS_OTEL_COLLECTOR_SERVICE_NAME,required"`
	OtelServiceVersion     string        `env:"METRICS_OTEL_COLLECTOR_SERVICE_VERSION,required"`
	OtelServiceEnvironment string        `env:"METRICS_OTEL_COLLECTOR_ENVIRONMENT,required"`
	OtelInterval           time.Duration `env:"METRICS_OTEL_COLLECTOR_INTERVAL,required"`
	OtelEndpoint           string        `env:"METRICS_OTEL_COLLECTOR_ENDPOINT,required"`
}

type metricsConfig struct {
	raw metricsEnvConfig
}

func NewMetricsConfig() (*metricsConfig, error) {
	var raw metricsEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &metricsConfig{raw: raw}, nil
}

func (cfg *metricsConfig) CollectorEndpoint() string {
	return cfg.raw.OtelEndpoint
}

func (cfg *metricsConfig) CollectorInterval() time.Duration {
	return cfg.raw.OtelInterval
}

func (cfg *metricsConfig) CollectorServiceName() string {
	return cfg.raw.OtelServiceName
}

func (cfg *metricsConfig) CollectorServiceVersion() string {
	return cfg.raw.OtelServiceVersion
}

func (cfg *metricsConfig) CollectorEnvironment() string {
	return cfg.raw.OtelServiceEnvironment
}
