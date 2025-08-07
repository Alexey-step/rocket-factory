package config

import (
	"os"

	"github.com/joho/godotenv"

	"github.com/Alexey-step/rocket-factory/order/internal/config/env"
)

var appConfig *config

type config struct {
	Logger    LoggerConfig
	OrderHTTP OrderHTTPConfig
	Inventory InventoryGRPCConfig
	Payment   PaymentGRPCConfig
	Postgres  PostgresConfig
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

	orderHTTPCfg, err := env.NewOrderHTTPConfig()
	if err != nil {
		return err
	}

	inventoryGRPCCfg, err := env.NewInventoryGRPCConfig()
	if err != nil {
		return err
	}

	paymentGRPCCfg, err := env.NewPaymentGRPCConfig()
	if err != nil {
		return err
	}

	postgresCfg, err := env.NewPostgresConfig()
	if err != nil {
		return err
	}

	appConfig = &config{
		Logger:    loggerCfg,
		OrderHTTP: orderHTTPCfg,
		Inventory: inventoryGRPCCfg,
		Payment:   paymentGRPCCfg,
		Postgres:  postgresCfg,
	}

	return nil
}

func AppConfig() *config {
	return appConfig
}
