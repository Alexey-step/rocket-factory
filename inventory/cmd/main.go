package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	inventoryV1API "github.com/Alexey-step/rocket-factory/inventory/internal/api/inventory/v1"
	interceptor "github.com/Alexey-step/rocket-factory/inventory/internal/interceptor"
	inventoryRepository "github.com/Alexey-step/rocket-factory/inventory/internal/repository/part"
	inventoryService "github.com/Alexey-step/rocket-factory/inventory/internal/service/part"
	inventoryV1 "github.com/Alexey-step/rocket-factory/shared/pkg/proto/inventory/v1"
)

const grpcPort = 50051

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Printf("failed to listen: %v\n", err)
		// return
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

	repo := inventoryRepository.NewRepository()
	service := inventoryService.NewService(repo)
	api := inventoryV1API.NewAPI(service)

	inventoryV1.RegisterInventoryServiceServer(s, api)

	reflection.Register(s)

	go func() {
		log.Printf("gRPC сервер запущен на порту %d\n", grpcPort)
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
