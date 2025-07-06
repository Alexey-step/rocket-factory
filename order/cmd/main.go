package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	customMiddleware "github.com/Alexey-step/rocket-factory/order/internal/middleware"
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
		return
	}

	inventoryClient := inventory_v1.NewInventoryServiceClient(inventoryConn)

	// Подключение к gRPC Payment-сервису
	paymentConn, err := grpc.NewClient(
		grpcPayment,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Printf("failed to connect to payment: %v\n", err)
	}

	paymentClient := payment_v1.NewPaymentServiceClient(paymentConn)

	// Создаем хранилище для данных о заказах
	storage := NewOrderStorage()
	service := NewOrderService(storage, inventoryClient, paymentClient)

	// Создаем обработчик API заказов деталей
	orderHandler := NewOrderHandler(service)

	// Создаем OpenAPI сервер
	orderServer, err := orderV1.NewServer(orderHandler)
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

type OrderStorage struct {
	mu     sync.RWMutex
	orders map[string]*orderV1.OrderDto
}

func NewOrderStorage() *OrderStorage {
	return &OrderStorage{
		orders: make(map[string]*orderV1.OrderDto),
	}
}

func NewOrderService(
	orderStorage *OrderStorage,
	inventoryClient inventory_v1.InventoryServiceClient,
	paymentClient payment_v1.PaymentServiceClient,
) *OrderService {
	return &OrderService{
		storage:          orderStorage,
		inventoryService: inventoryClient,
		paymentService:   paymentClient,
	}
}

type OrderService struct {
	storage          *OrderStorage
	paymentService   payment_v1.PaymentServiceClient
	inventoryService inventory_v1.InventoryServiceClient
}

// OrderHandler реализует интерфейс orderV1.Handler для обработки запросов к API заказа деталей
type OrderHandler struct {
	service *OrderService
}

// NewOrderHandler создает новый обработчик запросов к API заказа деталей
func NewOrderHandler(service *OrderService) *OrderHandler {
	return &OrderHandler{
		service: service,
	}
}

// GetOrder возвращает информацию о заказе
func (s *OrderStorage) GetOrder(uuid string) *orderV1.OrderDto {
	s.mu.RLock()
	defer s.mu.RUnlock()

	order, ok := s.orders[uuid]
	if !ok {
		return nil
	}

	return order
}

func (h *OrderHandler) CreateOrder(ctx context.Context, req *orderV1.CreateOrderRequest) (orderV1.CreateOrderRes, error) {
	parseUUIDs := uuidsToStrings(req.GetPartUuids())

	filter := &inventory_v1.PartsFilter{
		Uuids: parseUUIDs,
	}

	out, err := h.service.inventoryService.ListParts(ctx, &inventory_v1.ListPartsRequest{
		Filter: filter,
	})
	if err != nil {
		return &orderV1.NotFoundError{
			Code:    http.StatusNotFound,
			Message: "failed to retrieve parts: " + err.Error(),
		}, nil
	}

	parts := out.GetParts()

	if len(parts) != len(parseUUIDs) {
		return &orderV1.BadRequestError{
			Code:    http.StatusBadRequest,
			Message: "one or more parts not found",
		}, nil
	}

	var totalPrice float64
	for _, part := range parts {
		totalPrice += part.GetPrice()
	}

	orderUUID := uuid.New()

	order := &orderV1.OrderDto{
		OrderUUID:  orderUUID,
		UserUUID:   req.GetUserUUID(),
		PartUuids:  req.GetPartUuids(),
		TotalPrice: totalPrice,
		Status:     orderV1.OrderStatusPENDINGPAYMENT,
	}

	h.service.storage.mu.Lock()
	h.service.storage.orders[orderUUID.String()] = order
	h.service.storage.mu.Unlock()

	return &orderV1.CreateOrderResponse{
		OrderUUID:  orderUUID,
		TotalPrice: totalPrice,
	}, nil
}

func (h *OrderHandler) GetOrder(_ context.Context, params orderV1.GetOrderParams) (orderV1.GetOrderRes, error) {
	order := h.service.storage.GetOrder(params.OrderUUID.String())

	if order == nil {
		return &orderV1.NotFoundError{
			Code:    404,
			Message: "Order by this UUID`" + params.OrderUUID.String() + "` not found",
		}, nil
	}

	return &orderV1.GetOrderResponse{
		Data: *order,
	}, nil
}

func (h *OrderHandler) CancelOrder(_ context.Context, params orderV1.CancelOrderParams) (orderV1.CancelOrderRes, error) {
	order := h.service.storage.GetOrder(params.OrderUUID.String())

	if order == nil {
		return &orderV1.NotFoundError{
			Code:    404,
			Message: "Order by this UUID`" + params.OrderUUID.String() + "` not found",
		}, nil
	}

	switch order.Status {
	case orderV1.OrderStatusPAID:
		return &orderV1.ConflictError{
			Code:    409,
			Message: "Заказ уже оплачен и не может быть отменён",
		}, nil
	case orderV1.OrderStatusCANCELLED:
		return &orderV1.ConflictError{
			Code:    409,
			Message: "Заказ уже отменён",
		}, nil
	case orderV1.OrderStatusPENDINGPAYMENT:
		order.Status = orderV1.OrderStatusCANCELLED
		return nil, nil
	default:
		return &orderV1.InternalServerError{
			Code:    500,
			Message: "Неподдерживаемый статус заказа",
		}, nil
	}
}

func (h *OrderHandler) PayOrder(ctx context.Context, req *orderV1.PayOrderRequest, params orderV1.PayOrderParams) (orderV1.PayOrderRes, error) {
	order := h.service.storage.GetOrder(params.OrderUUID.String())

	if order == nil {
		return &orderV1.NotFoundError{
			Code:    http.StatusNotFound,
			Message: "Order by this UUID`" + params.OrderUUID.String() + "` not found",
		}, nil
	}

	switch order.Status {
	case orderV1.OrderStatusPAID:
		return &orderV1.ConflictError{
			Code:    http.StatusConflict,
			Message: "Заказ уже оплачен",
		}, nil
	case orderV1.OrderStatusCANCELLED:
		return &orderV1.ConflictError{
			Code:    http.StatusConflict,
			Message: "Заказ отменён и не может быть оплачен",
		}, nil
	case orderV1.OrderStatusPENDINGPAYMENT:
	default:
		return &orderV1.InternalServerError{
			Code:    http.StatusInternalServerError,
			Message: "Неподдерживаемый статус заказа",
		}, nil
	}

	paymentMethod := mapOrderToPaymentMethod(req.GetPaymentMethod())

	out, err := h.service.paymentService.PayOrder(ctx, &payment_v1.PayOrderRequest{
		OrderUuid:     params.OrderUUID.String(),
		PaymentMethod: paymentMethod,
		UserUuid:      uuid.New().String(),
	})
	if err != nil {
		return &orderV1.InternalServerError{
			Code:    http.StatusInternalServerError,
			Message: "Ошибка платежа: " + err.Error(),
		}, nil
	}

	parsedUUID, err := uuid.Parse(out.TransactionUuid)
	if err != nil {
		return &orderV1.InternalServerError{
			Code:    http.StatusInternalServerError,
			Message: "Некорректный UUID от платёжного сервиса",
		}, nil
	}

	order.Status = orderV1.OrderStatusPAID
	order.PaymentMethod = req.GetPaymentMethod()
	order.TransactionUUID = parsedUUID

	return &orderV1.PayOrderResponse{
		TransactionUUID: parsedUUID,
	}, nil
}

func (h *OrderHandler) NewError(_ context.Context, err error) *orderV1.GenericErrorStatusCode {
	return &orderV1.GenericErrorStatusCode{
		StatusCode: http.StatusInternalServerError,
		Response: orderV1.GenericError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		},
	}
}

func mapOrderToPaymentMethod(method orderV1.PaymentMethod) payment_v1.PaymentMethod {
	switch method {
	case orderV1.PaymentMethodCARD:
		return payment_v1.PaymentMethod_PAYMENT_METHOD_CARD
	case orderV1.PaymentMethodSBP:
		return payment_v1.PaymentMethod_PAYMENT_METHOD_SBP
	case orderV1.PaymentMethodCREDITCARD:
		return payment_v1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD
	case orderV1.PaymentMethodINVESTORMONEY:
		return payment_v1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY
	default:
		return payment_v1.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED
	}
}

func uuidsToStrings(uuids []uuid.UUID) []string {
	strs := make([]string, len(uuids))
	for i, u := range uuids {
		strs[i] = u.String()
	}
	return strs
}
