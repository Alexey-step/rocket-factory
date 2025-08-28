package app

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	inventoryV1API "github.com/Alexey-step/rocket-factory/inventory/internal/api/inventory/v1"
	grpcClient "github.com/Alexey-step/rocket-factory/inventory/internal/client/grpc"
	iamClient "github.com/Alexey-step/rocket-factory/inventory/internal/client/grpc/iam/v1"
	"github.com/Alexey-step/rocket-factory/inventory/internal/config"
	"github.com/Alexey-step/rocket-factory/inventory/internal/repository"
	inventoryRepository "github.com/Alexey-step/rocket-factory/inventory/internal/repository/part"
	"github.com/Alexey-step/rocket-factory/inventory/internal/service"
	inventoryService "github.com/Alexey-step/rocket-factory/inventory/internal/service/part"
	"github.com/Alexey-step/rocket-factory/platform/pkg/closer"
	authV1 "github.com/Alexey-step/rocket-factory/shared/pkg/proto/auth/v1"
	inventoryV1 "github.com/Alexey-step/rocket-factory/shared/pkg/proto/inventory/v1"
	userV1 "github.com/Alexey-step/rocket-factory/shared/pkg/proto/user/v1"
)

type diContainer struct {
	inventoryV1API      inventoryV1.InventoryServiceServer
	inventoryService    service.InventoryService
	inventoryRepository repository.InventoryRepository

	iamClient grpcClient.IamClient

	mongoDBClient *mongo.Client
	mongoDBHandle *mongo.Database
}

func NewDIContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) InventoryV1API(ctx context.Context) inventoryV1.InventoryServiceServer {
	if d.inventoryV1API == nil {
		d.inventoryV1API = inventoryV1API.NewAPI(d.PartService(ctx))
	}

	return d.inventoryV1API
}

func (d *diContainer) PartService(ctx context.Context) service.InventoryService {
	if d.inventoryService == nil {
		d.inventoryService = inventoryService.NewService(d.PartRepository(ctx))
	}

	return d.inventoryService
}

func (d *diContainer) PartRepository(ctx context.Context) repository.InventoryRepository {
	if d.inventoryRepository == nil {
		// тут глючит линтер, ругаясь что не передается контекст в NewRepository, пока кинул в игнор
		d.inventoryRepository = inventoryRepository.NewRepository(d.MongoDBHandle(ctx)) //nolint:contextcheck
	}

	return d.inventoryRepository
}

func (d *diContainer) MongoDBClient(ctx context.Context) *mongo.Client {
	if d.mongoDBClient == nil {
		client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.AppConfig().Mongo.URI()))
		if err != nil {
			panic(fmt.Sprintf("failed to connect to MongoDB: %s\n", err.Error()))
		}

		err = client.Ping(ctx, readpref.Primary())
		if err != nil {
			panic(fmt.Sprintf("failed to ping MongoDB: %s\n", err.Error()))
		}

		closer.AddNamed("MongoDB client", func(ctx context.Context) error {
			return client.Disconnect(ctx)
		})

		d.mongoDBClient = client
	}
	return d.mongoDBClient
}

func (d *diContainer) MongoDBHandle(ctx context.Context) *mongo.Database {
	if d.mongoDBHandle == nil {
		d.mongoDBHandle = d.MongoDBClient(ctx).Database(config.AppConfig().Mongo.DatabaseName())
	}

	return d.mongoDBHandle
}

func (d *diContainer) IamClient(_ context.Context) grpcClient.IamClient {
	if d.iamClient == nil {
		iamConn, err := grpc.NewClient(
			config.AppConfig().Iam.Address(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
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
