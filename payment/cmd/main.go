package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	interceptor "github.com/Alexey-step/rocket-factory/payment/internal/interceptor"
	paymentV1 "github.com/Alexey-step/rocket-factory/shared/pkg/proto/payment/v1"
)

const grpcPort = 50052

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Printf("failed to listen: %v\n", err)
		return
	}

	defer func() {
		if cerr := lis.Close(); cerr != nil {
			log.Printf("failed to close listener: %v\n", cerr)
		}
	}()

	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			grpc.UnaryServerInterceptor(interceptor.LoggerInterceptor()),
		),
	)

	service := &paymentService{}

	paymentV1.RegisterPaymentServiceServer(s, service)

	reflection.Register(s)

	go func() {
		log.Printf("starting gRPC server on port %d\n", grpcPort)
		err := s.Serve(lis)
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

type paymentService struct {
	paymentV1.UnimplementedPaymentServiceServer

	mu sync.RWMutex
}

func (s *paymentService) PayOrder(_ context.Context, _ *paymentV1.PayOrderRequest) (*paymentV1.PayOrderResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	UUID := uuid.New().String()

	log.Printf("✅Оплата прошла успешно, transaction_uuid: %v\n", UUID)

	return &paymentV1.PayOrderResponse{
		TransactionUuid: UUID,
	}, nil
}
