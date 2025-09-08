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

type OrderAssembledProducerConfig interface {
	TopicName() string
	Config() *sarama.Config
}

type OrderPaidConsumerConfig interface {
	Topic() string
	GroupID() string
	Config() *sarama.Config
}

type KafkaConfig interface {
	Brokers() []string
}

type MetricsConfig interface {
	CollectorEndpoint() string
	CollectorInterval() time.Duration
	CollectorServiceName() string
	CollectorServiceVersion() string
	CollectorEnvironment() string
}
