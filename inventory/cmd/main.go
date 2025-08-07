package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	inventoryV1API "github.com/Alexey-step/rocket-factory/inventory/internal/api/inventory/v1"
	interceptor "github.com/Alexey-step/rocket-factory/inventory/internal/interceptor"
	inventoryRepository "github.com/Alexey-step/rocket-factory/inventory/internal/repository/part"
	inventoryService "github.com/Alexey-step/rocket-factory/inventory/internal/service/part"
	inventoryV1 "github.com/Alexey-step/rocket-factory/shared/pkg/proto/inventory/v1"
)

const grpcAddr = "localhost:50051"

func main() {
	ctx := context.Background()

	err := godotenv.Load(".env")
	if err != nil {
		log.Printf("failed to load from .env file: %v\n", err)
		return
	}

	dbURI := os.Getenv("MONGO_URI")
	if dbURI == "" {
		log.Println("MONGO_URI environment variable is not set")
		return
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dbURI))
	if err != nil {
		log.Printf("failed to connect to MongoDB: %v\n", err)
		return
	}

	defer func() {
		if cerr := client.Disconnect(ctx); cerr != nil {
			log.Printf("failed to disconnect from MongoDB: %v\n", cerr)
		} else {
			log.Println("✅ Disconnected from MongoDB")
		}
	}()

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Printf("failed to ping MongoDB: %v\n", err)
		return
	} else {
		log.Println("✅ Connected to MongoDB")
	}

	db := client.Database("inventory-service")

	lis, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.Printf("failed to listen: %v\n", err)
		return
	}

	defer func() {
		if err = lis.Close(); err != nil {
			log.Printf("failed to close listener: %v\n", err)
		}
	}()

	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			interceptor.LoggerInterceptor(),
			interceptor.Validate(),
		),
	)

	repo := inventoryRepository.NewRepository(db)
	service := inventoryService.NewService(repo)
	api := inventoryV1API.NewAPI(service)

	inventoryV1.RegisterInventoryServiceServer(s, api)

	reflection.Register(s)

	go func() {
		log.Printf("🚀 gRPC InventoryService server listening on %s\n", grpcAddr)
		err = s.Serve(lis)
		if err != nil {
			log.Printf("failed to serve: %v\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("🛑Shutting down gRPC server...")
	s.GracefulStop()
	log.Println("✅ Server stopped")
}
