package app

import (
	"context"
	"fmt"
	"log"

	"github.com/IBM/sarama"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	v1 "github.com/Alexey-step/rocket-factory/order/internal/api/order/v1"
	grpcClient "github.com/Alexey-step/rocket-factory/order/internal/client/grpc"
	iamClient "github.com/Alexey-step/rocket-factory/order/internal/client/grpc/iam/v1"
	inventoryClient "github.com/Alexey-step/rocket-factory/order/internal/client/grpc/inventory/v1"
	paymentClient "github.com/Alexey-step/rocket-factory/order/internal/client/grpc/payment/v1"
	"github.com/Alexey-step/rocket-factory/order/internal/config"
	kafkaConverter "github.com/Alexey-step/rocket-factory/order/internal/converter/kafka"
	"github.com/Alexey-step/rocket-factory/order/internal/converter/kafka/decoder"
	"github.com/Alexey-step/rocket-factory/order/internal/repository"
	orderRepository "github.com/Alexey-step/rocket-factory/order/internal/repository/order"
	"github.com/Alexey-step/rocket-factory/order/internal/service"
	orderConsumer "github.com/Alexey-step/rocket-factory/order/internal/service/consumer/order_consumer"
	orderService "github.com/Alexey-step/rocket-factory/order/internal/service/order"
	orderProducer "github.com/Alexey-step/rocket-factory/order/internal/service/producer/order_producer"
	"github.com/Alexey-step/rocket-factory/platform/pkg/closer"
	wrappedKafka "github.com/Alexey-step/rocket-factory/platform/pkg/kafka"
	wrappedKafkaConsumer "github.com/Alexey-step/rocket-factory/platform/pkg/kafka/consumer"
	wrappedKafkaProducer "github.com/Alexey-step/rocket-factory/platform/pkg/kafka/producer"
	"github.com/Alexey-step/rocket-factory/platform/pkg/logger"
	kafkaMiddleware "github.com/Alexey-step/rocket-factory/platform/pkg/middleware/kafka"
	"github.com/Alexey-step/rocket-factory/platform/pkg/migrator"
	pgMigrator "github.com/Alexey-step/rocket-factory/platform/pkg/migrator/pg"
	"github.com/Alexey-step/rocket-factory/platform/pkg/tracing"
	orderV1 "github.com/Alexey-step/rocket-factory/shared/pkg/openapi/order/v1"
	authV1 "github.com/Alexey-step/rocket-factory/shared/pkg/proto/auth/v1"
	inventory_v1 "github.com/Alexey-step/rocket-factory/shared/pkg/proto/inventory/v1"
	paymentV1 "github.com/Alexey-step/rocket-factory/shared/pkg/proto/payment/v1"
	userV1 "github.com/Alexey-step/rocket-factory/shared/pkg/proto/user/v1"
)

type diContainer struct {
	orderV1API      orderV1.Handler
	orderService    service.OrderService
	orderRepository repository.OrderRepository

	inventoryClient grpcClient.InventoryClient
	paymentClient   grpcClient.PaymentClient
	iamClient       grpcClient.IamClient

	postgresDB *pgxpool.Pool
	migrator   migrator.Migrator

	orderProducerService service.OrderProducerService
	orderConsumerService service.OrderConsumerService

	consumerGroup sarama.ConsumerGroup
	syncProducer  sarama.SyncProducer

	orderPaidProducer      wrappedKafka.Producer
	orderAssembledConsumer wrappedKafka.Consumer
	orderAssembledDecoder  kafkaConverter.OrderAssembledDecoder
}

func NewDIContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) OrderV1API(ctx context.Context) orderV1.Handler {
	if d.orderV1API == nil {
		d.orderV1API = v1.NewAPI(d.OrderService(ctx))
	}

	return d.orderV1API
}

func (d *diContainer) OrderService(ctx context.Context) service.OrderService {
	if d.orderService == nil {
		d.orderService = orderService.NewService(
			d.OrderRepository(ctx),
			d.InventoryClient(ctx),
			d.PaymentClient(ctx),
			d.OrderProducerService(ctx),
		)
	}
	return d.orderService
}

func (d *diContainer) InventoryClient(_ context.Context) grpcClient.InventoryClient {
	if d.inventoryClient == nil {
		inventoryConn, err := grpc.NewClient(
			config.AppConfig().Inventory.Address(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithUnaryInterceptor(tracing.UnaryClientInterceptor("order-service")),
		)
		if err != nil {
			log.Printf("failed to connect to invertory: %v\n", err)
		}

		closer.AddNamed("Inventory client", func(ctx context.Context) error {
			return inventoryConn.Close()
		})

		d.inventoryClient = inventoryClient.NewClient(inventory_v1.NewInventoryServiceClient(inventoryConn))
	}
	return d.inventoryClient
}

func (d *diContainer) PaymentClient(_ context.Context) grpcClient.PaymentClient {
	if d.paymentClient == nil {
		paymentConn, err := grpc.NewClient(
			config.AppConfig().Payment.Address(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithUnaryInterceptor(tracing.UnaryClientInterceptor("order-service")),
		)
		if err != nil {
			log.Printf("failed to connect to payment: %v\n", err)
		}

		closer.AddNamed("Payment client", func(ctx context.Context) error {
			return paymentConn.Close()
		})

		d.paymentClient = paymentClient.NewClient(paymentV1.NewPaymentServiceClient(paymentConn))
	}
	return d.paymentClient
}

func (d *diContainer) IamClient(_ context.Context) grpcClient.IamClient {
	if d.iamClient == nil {
		iamConn, err := grpc.NewClient(
			config.AppConfig().Iam.Address(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithUnaryInterceptor(tracing.UnaryClientInterceptor("order-service")),
		)
		if err != nil {
			log.Printf("failed to connect to iam: %v\n", err)
		}

		closer.AddNamed("Iam client", func(ctx context.Context) error {
			return iamConn.Close()
		})

		d.iamClient = iamClient.NewClient(
			authV1.NewAuthServiceClient(iamConn),
			userV1.NewUserServiceClient(iamConn),
		)
	}
	return d.iamClient
}

func (d *diContainer) OrderRepository(ctx context.Context) repository.OrderRepository {
	if d.orderRepository == nil {
		d.orderRepository = orderRepository.NewOrderRepository(d.PostgresDB(ctx))
	}
	return d.orderRepository
}

func (d *diContainer) PostgresDB(ctx context.Context) *pgxpool.Pool {
	if d.postgresDB == nil {
		pool, err := pgxpool.New(ctx, config.AppConfig().Postgres.URI())
		if err != nil {
			panic(fmt.Sprintf("failed to connect to Postgres: %s\n", err.Error()))
		}

		// Проверяем соединение с базой данных
		err = pool.Ping(ctx)
		if err != nil {
			panic(fmt.Sprintf("failed to ping Postgres: %s\n", err.Error()))
		}

		closer.AddNamed("PostgresDB", func(ctx context.Context) error {
			pool.Close()
			return nil
		})

		d.postgresDB = pool
	}

	return d.postgresDB
}

func (d *diContainer) Migrator(_ context.Context) migrator.Migrator {
	if d.migrator == nil {
		cfg, err := pgxpool.ParseConfig(config.AppConfig().Postgres.URI())
		if err != nil {
			panic(fmt.Sprintf("failed to parse Postgres config: %s\n", err.Error()))
		}

		migrationsDir := config.AppConfig().OrderHTTP.MigrationsDir()
		d.migrator = pgMigrator.NewMigrator(stdlib.OpenDB(*cfg.ConnConfig), migrationsDir)
	}

	return d.migrator
}

func (d *diContainer) OrderConsumerService(ctx context.Context) service.OrderConsumerService {
	if d.orderConsumerService == nil {
		d.orderConsumerService = orderConsumer.NewService(
			d.OrderAssembledDecoder(),
			d.OrderAssembledConsumer(),
			d.OrderService(ctx),
		)
	}
	return d.orderConsumerService
}

func (d *diContainer) OrderProducerService(_ context.Context) service.OrderProducerService {
	if d.orderProducerService == nil {
		d.orderProducerService = orderProducer.NewService(
			d.OrderPaidProducer(),
		)
	}
	return d.orderProducerService
}

func (d *diContainer) ConsumerGroup() sarama.ConsumerGroup {
	if d.consumerGroup == nil {
		group, err := sarama.NewConsumerGroup(
			config.AppConfig().Kafka.Brokers(),
			config.AppConfig().OrderAssembledConsumer.GroupID(),
			config.AppConfig().OrderAssembledConsumer.Config(),
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

func (d *diContainer) SyncProducer() sarama.SyncProducer {
	if d.syncProducer == nil {
		p, err := sarama.NewSyncProducer(
			config.AppConfig().Kafka.Brokers(),
			config.AppConfig().OrderPaidProducer.Config(),
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

func (d *diContainer) OrderPaidProducer() wrappedKafka.Producer {
	if d.orderPaidProducer == nil {
		d.orderPaidProducer = wrappedKafkaProducer.NewProducer(
			d.SyncProducer(),
			config.AppConfig().OrderPaidProducer.TopicName(),
			logger.Logger(),
		)
	}
	return d.orderPaidProducer
}
