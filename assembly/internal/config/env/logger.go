package env

import (
	"github.com/caarlos0/env/v11"
)

type loggerEnvConfig struct {
	Level                  string `env:"LOGGER_LEVEL,required"`
	AsJson                 bool   `env:"LOGGER_AS_JSON,required"`
	OtelEnabled            bool   `env:"LOGGER_OTEL_COLLECTOR_ENABLED,required"`
	OtelServiceName        string `env:"LOGGER_OTEL_COLLECTOR_SERVICE_NAME,required"`
	OtelServiceEnvironment string `env:"LOGGER_OTEL_COLLECTOR_ENVIRONMENT,required"`
	OtelEndpoint           string `env:"LOGGER_OTEL_COLLECTOR_ENDPOINT,required"`
}

type loggerConfig struct {
	raw loggerEnvConfig
}

func NewLoggerConfig() (*loggerConfig, error) {
	var raw loggerEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &loggerConfig{raw: raw}, nil
}

func (cfg *loggerConfig) Level() string {
	return cfg.raw.Level
}

func (cfg *loggerConfig) AsJson() bool {
	return cfg.raw.AsJson
}

func (cfg *loggerConfig) OtelEnabled() bool {
	return cfg.raw.OtelEnabled
}

func (cfg *loggerConfig) OtelServiceName() string {
	return cfg.raw.OtelServiceName
}

func (cfg *loggerConfig) OtelServiceEnvironment() string {
	return cfg.raw.OtelServiceEnvironment
}

func (cfg *loggerConfig) OtelEndpoint() string {
	return cfg.raw.OtelEndpoint
}
