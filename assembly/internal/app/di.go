package app

import (
	"context"
	"fmt"

	"github.com/IBM/sarama"

	"github.com/Alexey-step/rocket-factory/assembly/internal/config"
	kafkaConverter "github.com/Alexey-step/rocket-factory/assembly/internal/converter/kafka"
	"github.com/Alexey-step/rocket-factory/assembly/internal/converter/kafka/decoder"
	"github.com/Alexey-step/rocket-factory/assembly/internal/service"
	orderConsumer "github.com/Alexey-step/rocket-factory/assembly/internal/service/consumer/order_consumer"
	orderProducer "github.com/Alexey-step/rocket-factory/assembly/internal/service/producer/order_producer"
	"github.com/Alexey-step/rocket-factory/platform/pkg/closer"
	wrappedKafka "github.com/Alexey-step/rocket-factory/platform/pkg/kafka"
	wrappedKafkaConsumer "github.com/Alexey-step/rocket-factory/platform/pkg/kafka/consumer"
	wrappedKafkaProducer "github.com/Alexey-step/rocket-factory/platform/pkg/kafka/producer"
	"github.com/Alexey-step/rocket-factory/platform/pkg/logger"
	kafkaMiddleware "github.com/Alexey-step/rocket-factory/platform/pkg/middleware/kafka"
)

type diContainer struct {
	assemblyService      service.AssemblyService
	orderProducerService service.OrderProducerService
	orderConsumerService service.OrderConsumerService

	consumerGroup sarama.ConsumerGroup
	syncProducer  sarama.SyncProducer

	orderAssembledProducer wrappedKafka.Producer
	orderPaidConsumer      wrappedKafka.Consumer
	orderPaidDecoder       kafkaConverter.OrderPaidDecoder
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) AssemblyService() service.AssemblyService {
	if d.assemblyService == nil {
		d.assemblyService = service.NewService(
			d.OrderProducerService())
	}
	return d.assemblyService
}

func (d *diContainer) OrderConsumerService() service.OrderConsumerService {
	if d.orderConsumerService == nil {
		d.orderConsumerService = orderConsumer.NewService(
			d.OrderPaidConsumer(),
			d.OrderPaidDecoder(),
			d.OrderProducerService(),
		)
	}
	return d.orderConsumerService
}

func (d *diContainer) OrderProducerService() service.OrderProducerService {
	if d.orderProducerService == nil {
		d.orderProducerService = orderProducer.NewService(
			d.OrderAssembledProducer(),
		)
	}
	return d.orderProducerService
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

func (d *diContainer) SyncProducer() sarama.SyncProducer {
	if d.syncProducer == nil {
		p, err := sarama.NewSyncProducer(
			config.AppConfig().Kafka.Brokers(),
			config.AppConfig().OrderAssembledProducer.Config(),
		)
		if err != nil {
			panic(fmt.Sprintf("failed to create sync producer: %s\n", err.Error()))
		}
		closer.AddNamed("Kafka sync producer", func(ctx context.Context) error {
			return p.Close()
		})

		d.syncProducer = p
	}

	return d.syncProducer
}

func (d *diContainer) OrderAssembledProducer() wrappedKafka.Producer {
	if d.orderAssembledProducer == nil {
		d.orderAssembledProducer = wrappedKafkaProducer.NewProducer(
			d.SyncProducer(),
			config.AppConfig().OrderAssembledProducer.TopicName(),
			logger.Logger(),
		)
	}
	return d.orderAssembledProducer
}
