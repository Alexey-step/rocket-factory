package config

type LoggerConfig interface {
	Level() string
	AsJson() bool
	OtelEnabled() bool
	OtelServiceName() string
	OtelServiceEnvironment() string
	OtelEndpoint() string
}

type PaymentGRPCConfig interface {
	Address() string
}

type TracingConfig interface {
	CollectorEndpoint() string
	ServiceName() string
	Environment() string
	ServiceVersion() string
}
