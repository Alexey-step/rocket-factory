package app

import (
	"context"
	"fmt"

	"github.com/IBM/sarama"
	"github.com/go-telegram/bot"

	httpClient "github.com/Alexey-step/rocket-factory/notification/internal/client/http"
	telegramClient "github.com/Alexey-step/rocket-factory/notification/internal/client/http/telegram"
	"github.com/Alexey-step/rocket-factory/notification/internal/config"
	kafkaConverter "github.com/Alexey-step/rocket-factory/notification/internal/converter/kafka"
	"github.com/Alexey-step/rocket-factory/notification/internal/converter/kafka/decoder"
	"github.com/Alexey-step/rocket-factory/notification/internal/service"
	orderAssembledConsumerService "github.com/Alexey-step/rocket-factory/notification/internal/service/consumer/order_assembled_consumer"
	orderPaidConsumerService "github.com/Alexey-step/rocket-factory/notification/internal/service/consumer/order_paid_consumer"
	telegramService "github.com/Alexey-step/rocket-factory/notification/internal/service/telegram"
	"github.com/Alexey-step/rocket-factory/platform/pkg/closer"
	wrappedKafka "github.com/Alexey-step/rocket-factory/platform/pkg/kafka"
	wrappedKafkaConsumer "github.com/Alexey-step/rocket-factory/platform/pkg/kafka/consumer"
	"github.com/Alexey-step/rocket-factory/platform/pkg/logger"
	kafkaMiddleware "github.com/Alexey-step/rocket-factory/platform/pkg/middleware/kafka"
)

type diContainer struct {
	consumerGroup sarama.ConsumerGroup

	orderAssembledConsumerService service.OrderAssembledConsumerService
	orderPaidConsumerService      service.OrderPaidConsumerService

	orderAssembledConsumer wrappedKafka.Consumer
	orderAssembledDecoder  kafkaConverter.OrderAssembledDecoder
	orderPaidConsumer      wrappedKafka.Consumer
	orderPaidDecoder       kafkaConverter.OrderPaidDecoder

	telegramService service.TelegramService
	telegramClient  httpClient.TelegramClient
	telegramBot     *bot.Bot
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) TelegramService(ctx context.Context) service.TelegramService {
	if d.telegramService == nil {
		d.telegramService = telegramService.NewService(
			d.TelegramClient(ctx),
		)
	}
	return d.telegramService
}

func (d *diContainer) TelegramClient(ctx context.Context) httpClient.TelegramClient {
	if d.telegramClient == nil {
		d.telegramClient = telegramClient.NewClient(d.TelegramBot(ctx))
	}
	return d.telegramClient
}

func (d *diContainer) TelegramBot(_ context.Context) *bot.Bot {
	if d.telegramBot == nil {
		b, err := bot.New(config.AppConfig().TelegramBot.Token())
		if err != nil {
			panic("failed to create Telegram bot: " + err.Error())
		}

		d.telegramBot = b
	}
	return d.telegramBot
}

func (d *diContainer) OrderAssembledConsumerService(ctx context.Context) service.OrderAssembledConsumerService {
	if d.orderAssembledConsumerService == nil {
		d.orderAssembledConsumerService = orderAssembledConsumerService.NewService(
			d.OrderAssembledConsumer(),
			d.OrderAssembledDecoder(),
			d.TelegramService(ctx),
		)
	}
	return d.orderAssembledConsumerService
}

func (d *diContainer) ConsumerGroup() sarama.ConsumerGroup {
	if d.consumerGroup == nil {
		group, err := sarama.NewConsumerGroup(
			config.AppConfig().Kafka.Brokers(),
			config.AppConfig().OrderPaidConsumer.GroupID(),
			config.AppConfig().OrderPaidConsumer.Config(),
		)
		if err != nil {
			panic(fmt.Sprintf("failed to create consumer group: %s\n", err.Error()))
		}
		closer.AddNamed("Kafka consumer group", func(ctx context.Context) error {
			return group.Close()
		})

		d.consumerGroup = group
	}

	return d.consumerGroup
}

func (d *diContainer) OrderAssembledConsumer() wrappedKafka.Consumer {
	if d.orderAssembledConsumer == nil {
		d.orderAssembledConsumer = wrappedKafkaConsumer.NewConsumer(
			logger.Logger(),
			d.ConsumerGroup(),
			[]string{
				config.AppConfig().OrderAssembledConsumer.Topic(),
			},
			kafkaMiddleware.Logging(logger.Logger()),
		)
	}
	return d.orderAssembledConsumer
}

func (d *diContainer) OrderAssembledDecoder() kafkaConverter.OrderAssembledDecoder {
	if d.orderAssembledDecoder == nil {
		d.orderAssembledDecoder = decoder.NewOrderAssembledDecoder()
	}
	return d.orderAssembledDecoder
}

func (d *diContainer) OrderPaidConsumerService(ctx context.Context) service.OrderPaidConsumerService {
	if d.orderPaidConsumerService == nil {
		d.orderPaidConsumerService = orderPaidConsumerService.NewService(
			d.OrderPaidConsumer(),
			d.OrderPaidDecoder(),
			d.TelegramService(ctx),
		)
	}
	return d.orderPaidConsumerService
}

func (d *diContainer) OrderPaidConsumer() wrappedKafka.Consumer {
	if d.orderPaidConsumer == nil {
		d.orderPaidConsumer = wrappedKafkaConsumer.NewConsumer(
			logger.Logger(),
			d.ConsumerGroup(),
			[]string{
				config.AppConfig().OrderPaidConsumer.Topic(),
			},
			kafkaMiddleware.Logging(logger.Logger()),
		)
	}
	return d.orderPaidConsumer
}

func (d *diContainer) OrderPaidDecoder() kafkaConverter.OrderPaidDecoder {
	if d.orderPaidDecoder == nil {
		d.orderPaidDecoder = decoder.NewOrderPaidDecoder()
	}
	return d.orderPaidDecoder
}
