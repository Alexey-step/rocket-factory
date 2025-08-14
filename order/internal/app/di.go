package app

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	v1 "github.com/Alexey-step/rocket-factory/order/internal/api/order/v1"
	grpcClient "github.com/Alexey-step/rocket-factory/order/internal/client/grpc"
	inventoryClient "github.com/Alexey-step/rocket-factory/order/internal/client/grpc/inventory/v1"
	paymentClient "github.com/Alexey-step/rocket-factory/order/internal/client/grpc/payment/v1"
	"github.com/Alexey-step/rocket-factory/order/internal/config"
	"github.com/Alexey-step/rocket-factory/order/internal/repository"
	orderRepository "github.com/Alexey-step/rocket-factory/order/internal/repository/order"
	"github.com/Alexey-step/rocket-factory/order/internal/service"
	orderService "github.com/Alexey-step/rocket-factory/order/internal/service/order"
	"github.com/Alexey-step/rocket-factory/platform/pkg/closer"
	"github.com/Alexey-step/rocket-factory/platform/pkg/migrator"
	pgMigrator "github.com/Alexey-step/rocket-factory/platform/pkg/migrator/pg"
	orderV1 "github.com/Alexey-step/rocket-factory/shared/pkg/openapi/order/v1"
	inventory_v1 "github.com/Alexey-step/rocket-factory/shared/pkg/proto/inventory/v1"
	paymentV1 "github.com/Alexey-step/rocket-factory/shared/pkg/proto/payment/v1"
)

type diContainer struct {
	orderV1API      orderV1.Handler
	orderService    service.OrderService
	orderRepository repository.OrderRepository

	inventoryClient grpcClient.InventoryClient
	paymentClient   grpcClient.PaymentClient

	postgresDB *pgxpool.Pool
	migrator   migrator.Migrator
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
		d.orderService = orderService.NewService(d.OrderRepository(ctx), d.InventoryClient(ctx), d.PaymentClient(ctx))
	}
	return d.orderService
}

func (d *diContainer) InventoryClient(_ context.Context) grpcClient.InventoryClient {
	if d.inventoryClient == nil {
		inventoryConn, err := grpc.NewClient(
			config.AppConfig().Inventory.Address(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
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
		)
		if err != nil {
			log.Printf("failed to connect to invertory: %v\n", err)
		}

		closer.AddNamed("Payment client", func(ctx context.Context) error {
			return paymentConn.Close()
		})

		d.paymentClient = paymentClient.NewClient(paymentV1.NewPaymentServiceClient(paymentConn))
	}
	return d.paymentClient
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
