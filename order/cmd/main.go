package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	orderV1Api "github.com/Alexey-step/rocket-factory/order/internal/api/order/v1"
	inventoryV1Client "github.com/Alexey-step/rocket-factory/order/internal/client/grpc/inventory/v1"
	paymentV1Client "github.com/Alexey-step/rocket-factory/order/internal/client/grpc/payment/v1"
	customMiddleware "github.com/Alexey-step/rocket-factory/order/internal/middleware"
	orderRepository "github.com/Alexey-step/rocket-factory/order/internal/repository/order"
	orderService "github.com/Alexey-step/rocket-factory/order/internal/service/order"
	orderV1 "github.com/Alexey-step/rocket-factory/shared/pkg/openapi/order/v1"
	inventory_v1 "github.com/Alexey-step/rocket-factory/shared/pkg/proto/inventory/v1"
	payment_v1 "github.com/Alexey-step/rocket-factory/shared/pkg/proto/payment/v1"
)

const (
	httpPort          = 8080
	readHeaderTimeout = 5 * time.Second
	shutdownTimeout   = 10 * time.Second

	grpcInventory = "localhost:50051"
	grpcPayment   = "localhost:50052"
)

func main() {
	// Подключение к gRPC Inventory-сервису
	inventoryConn, err := grpc.NewClient(
		grpcInventory,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Printf("failed to connect to invertory: %v\n", err)
		// return
	}

	inventoryClient := inventoryV1Client.NewClient(inventory_v1.NewInventoryServiceClient(inventoryConn))

	// Подключение к gRPC Payment-сервису
	paymentConn, err := grpc.NewClient(
		grpcPayment,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Printf("failed to connect to payment: %v\n", err)
	}

	paymentClient := paymentV1Client.NewClient(payment_v1.NewPaymentServiceClient(paymentConn))

	repo := orderRepository.NewOrderRepository()
	service := orderService.NewService(repo, inventoryClient, paymentClient)
	api := orderV1Api.NewAPI(service)

	// Создаем OpenAPI сервер
	orderServer, err := orderV1.NewServer(api)
	if err != nil {
		if err := inventoryConn.Close(); err != nil {
			log.Printf("failed to close inventory connection: %v\n", err)
		}
		if err := paymentConn.Close(); err != nil {
			log.Printf("failed to close payment connection: %v\n", err)
		}
		log.Fatalf("ошибка создания сервера OpenAPI: %v", err)
	}

	defer func() {
		if err := inventoryConn.Close(); err != nil {
			log.Printf("failed to close inventory connection: %v\n", err)
		}
	}()
	defer func() {
		if err := paymentConn.Close(); err != nil {
			log.Printf("failed to close payment connection: %v\n", err)
		}
	}()

	r := chi.NewRouter()

	// Добавляем middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(customMiddleware.RequestLogger)
	r.Use(middleware.Timeout(10 * time.Second))

	r.Mount("/", orderServer)

	server := &http.Server{
		Addr:        fmt.Sprintf(":%d", httpPort),
		Handler:     r,
		ReadTimeout: readHeaderTimeout, // Защита от Slowloris атак - тип DDoS-атаки, при которой
		// атакующий умышленно медленно отправляет HTTP-заголовки, удерживая соединения открытыми и истощая
		// пул доступных соединений на сервере. ReadHeaderTimeout принудительно закрывает соединение,
		// если клиент не успел отправить все заголовки за отведенное время.
	}

	// Запускаем сервер в отдельной горутине
	go func() {
		log.Printf("🚀 HTTP-сервер запущен на порту %d\n", httpPort)
		err = server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("❌ Ошибка запуска сервера: %v\n", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Завершение работы сервера...")

	// Создаем контекст с таймаутом для остановки сервера
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	err = server.Shutdown(ctx)
	if err != nil {
		log.Printf("❌ Ошибка при остановке сервера: %v\n", err)
	}

	log.Println("✅ Сервер остановлен")
}
