package config

import (
	"time"

	"github.com/IBM/sarama"
)

type LoggerConfig interface {
	Level() string
	AsJson() bool
	OtelEnabled() bool
	OtelServiceName() string
	OtelServiceEnvironment() string
	OtelEndpoint() string
}

type OrderHTTPConfig interface {
	Address() string
	MigrationsDir() string
}

type InventoryGRPCConfig interface {
	Address() string
}

type PaymentGRPCConfig interface {
	Address() string
}

type PostgresConfig interface {
	URI() string
	DatabaseName() string
}

type OrderAssembledConsumerConfig interface {
	Topic() string
	GroupID() string
	Config() *sarama.Config
}

type OrderPaidProducerConfig interface {
	TopicName() string
	Config() *sarama.Config
}

type KafkaConfig interface {
	Brokers() []string
}

type IamGRPCConfig interface {
	Address() string
}

type TracingConfig interface {
	CollectorEndpoint() string
	ServiceName() string
	Environment() string
	ServiceVersion() string
}

type MetricsConfig interface {
	CollectorEndpoint() string
	CollectorInterval() time.Duration
	CollectorServiceName() string
	CollectorServiceVersion() string
	CollectorEnvironment() string
}
