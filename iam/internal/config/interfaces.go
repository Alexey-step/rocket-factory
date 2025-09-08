package config

import "time"

type IamGRPCConfig interface {
	Address() string
	MigrationsDir() string
}

type LoggerConfig interface {
	Level() string
	AsJson() bool
	OtelEnabled() bool
	OtelServiceName() string
	OtelServiceEnvironment() string
	OtelEndpoint() string
}

type SessionConfig interface {
	TTL() time.Duration
}

type RedisConfig interface {
	Host() string
	Port() string
	ConnectionTimeout() time.Duration
	MaxIDLE() int
	IDLETimeout() time.Duration
	Address() string
}

type PostgresConfig interface {
	URI() string
	DatabaseName() string
}
