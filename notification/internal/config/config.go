package config

import (
	"os"

	"github.com/joho/godotenv"

	"github.com/Alexey-step/rocket-factory/notification/internal/config/env"
)

var appConfig *config

type config struct {
	Logger                 LoggerConfig
	Kafka                  KafkaConfig
	OrderAssembledConsumer OrderAssembledConsumerConfig
	OrderPaidConsumer      OrderPaidConsumerConfig
	TelegramBot            TelegramBotConfig
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

	kafkaCfg, err := env.NewKafkaConfig()
	if err != nil {
		return err
	}

	orderAssembledConsumerCfg, err := env.NewOrderAssembledConsumerConfig()
	if err != nil {
		return err
	}

	orderPaidConsumerCfg, err := env.NewOrderPaidConsumerConfig()
	if err != nil {
		return err
	}

	telegramBotConfig, err := env.NewTelegramBotConfig()
	if err != nil {
		return err
	}

	appConfig = &config{
		Logger:                 loggerCfg,
		Kafka:                  kafkaCfg,
		OrderAssembledConsumer: orderAssembledConsumerCfg,
		OrderPaidConsumer:      orderPaidConsumerCfg,
		TelegramBot:            telegramBotConfig,
	}

	return nil
}

func AppConfig() *config {
	return appConfig
}
