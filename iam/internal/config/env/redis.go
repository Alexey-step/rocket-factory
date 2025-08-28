package env

import (
	"net"
	"time"

	"github.com/caarlos0/env/v11"
)

type redisEnvConfig struct {
	Host              string        `env:"REDIS_HOST,required"`
	Port              string        `env:"REDIS_PORT,required"`
	ConnectionTimeout time.Duration `env:"REDIS_CONNECTION_TIMEOUT,required"`
	MaxIDLE           int           `env:"REDIS_MAX_IDLE,required"`
	IDLETimeout       time.Duration `env:"REDIS_IDLE_TIMEOUT,required"`
}

type redisConfig struct {
	raw redisEnvConfig
}

func NewRedisConfig() (*redisConfig, error) {
	var raw redisEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &redisConfig{raw: raw}, nil
}

func (cfg *redisConfig) Host() string {
	return cfg.raw.Host
}

func (cfg *redisConfig) Port() string {
	return cfg.raw.Port
}

func (cfg *redisConfig) ConnectionTimeout() time.Duration {
	return cfg.raw.ConnectionTimeout
}

func (cfg *redisConfig) MaxIDLE() int {
	return cfg.raw.MaxIDLE
}

func (cfg *redisConfig) IDLETimeout() time.Duration {
	return cfg.raw.IDLETimeout
}

func (cfg *redisConfig) Address() string {
	return net.JoinHostPort(cfg.raw.Host, cfg.raw.Port)
}
