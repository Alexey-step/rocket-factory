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

	paymentV1Api "github.com/Alexey-step/rocket-factory/payment/internal/api/payment/v1"
	interceptor "github.com/Alexey-step/rocket-factory/payment/internal/interceptor"
	paymentService "github.com/Alexey-step/rocket-factory/payment/internal/service/payment"
	paymentV1 "github.com/Alexey-step/rocket-factory/shared/pkg/proto/payment/v1"
)

const grpcPort = 50052

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
		),
	)

	service := paymentService.NewService()
	api := paymentV1Api.NewApi(service)

	paymentV1.RegisterPaymentServiceServer(s, api)

	reflection.Register(s)

	go func() {
		log.Printf("starting gRPC server on port %d\n", grpcPort)
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
