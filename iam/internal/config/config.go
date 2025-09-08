package config

import (
	"os"

	"github.com/joho/godotenv"

	"github.com/Alexey-step/rocket-factory/iam/internal/config/env"
)

var appConfig *config

type config struct {
	Logger   LoggerConfig
	IamGRPC  IamGRPCConfig
	Postgres PostgresConfig
	Session  SessionConfig
	Redis    RedisConfig
}

func Load(path ...string) error {
	err := godotenv.Load(path...)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	loggerCfg, err := env.NewLoggerConfig()
	if err != nil {
		return err
	}

	iamGRPCConfig, err := env.NewIamGRPCConfig()
	if err != nil {
		return err
	}

	postgresCfg, err := env.NewPostgresConfig()
	if err != nil {
		return err
	}

	sessionCfg, err := env.NewSessionConfig()
	if err != nil {
		return err
	}

	redisCfg, err := env.NewRedisConfig()
	if err != nil {
		return err
	}

	appConfig = &config{
		Logger:   loggerCfg,
		IamGRPC:  iamGRPCConfig,
		Postgres: postgresCfg,
		Session:  sessionCfg,
		Redis:    redisCfg,
	}

	return nil
}

func AppConfig() *config {
	return appConfig
}
