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

	Kafka                  KafkaConfig
	OrderPaidProducer      OrderPaidProducerConfig
	OrderAssembledConsumer OrderAssembledConsumerConfig
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

	kafkaCfg, err := env.NewKafkaConfig()
	if err != nil {
		return err
	}

	orderPaidProducerCfg, err := env.NewOrderPaidProducerConfig()
	if err != nil {
		return err
	}

	orderAssembledConsumerCfg, err := env.NewOrderAssembledConsumerConfig()
	if err != nil {
		return err
	}

	appConfig = &config{
		Logger:                 loggerCfg,
		OrderHTTP:              orderHTTPCfg,
		Inventory:              inventoryGRPCCfg,
		Payment:                paymentGRPCCfg,
		Postgres:               postgresCfg,
		Kafka:                  kafkaCfg,
		OrderPaidProducer:      orderPaidProducerCfg,
		OrderAssembledConsumer: orderAssembledConsumerCfg,
	}

	return nil
}

func AppConfig() *config {
	return appConfig
}
