package main

import (
	"context"
	"fmt"
	"log"
	"math"
	"net"
	"os"
	"os/signal"
	"slices"
	"sync"
	"syscall"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	interceptor "github.com/Alexey-step/rocket-factory/inventory/internal/interceptor"
	inventoryV1 "github.com/Alexey-step/rocket-factory/shared/pkg/proto/inventory/v1"
)

const grpcPort = 50051

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

	storage := NewInventoryStorageInMemory()
	service := &InventoryService{
		storage: storage,
	}

	inventoryV1.RegisterInventoryServiceServer(s, service)

	reflection.Register(s)

	go func() {
		log.Printf("gRPC сервер запущен на порту %d\n", grpcPort)
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

type InventoryStorage interface {
	Part(uuid string) (*inventoryV1.Part, error)
	Parts(filter *inventoryV1.PartsFilter) ([]*inventoryV1.Part, error)
}

type InventoryStorageInMemory struct {
	mu    sync.RWMutex
	parts map[string]*inventoryV1.Part
}

func NewInventoryStorageInMemory() *InventoryStorageInMemory {
	s := &InventoryStorageInMemory{
		parts: make(map[string]*inventoryV1.Part),
	}
	s.initParts()
	return s
}

func (s *InventoryStorageInMemory) Part(uuid string) (*inventoryV1.Part, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	part, ok := s.parts[uuid]
	if !ok {
		return nil, fmt.Errorf("part with UUID %s not found", uuid)
	}

	return part, nil
}

func (s *InventoryStorageInMemory) Parts(filter *inventoryV1.PartsFilter) ([]*inventoryV1.Part, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []*inventoryV1.Part
	for _, part := range s.parts {
		if matchesFilter(part, filter) {
			result = append(result, part)
		}
	}
	return result, nil
}

func (s *InventoryStorageInMemory) initParts() {
	for _, part := range generateParts() {
		s.parts[part.Uuid] = part
	}
}

type InventoryService struct {
	inventoryV1.UnimplementedInventoryServiceServer
	storage InventoryStorage
}

func (s *InventoryService) GetPart(_ context.Context, req *inventoryV1.GetPartRequest) (*inventoryV1.GetPartResponse, error) {
	part, err := s.storage.Part(req.GetUuid())
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "error: %v", err)
	}

	return &inventoryV1.GetPartResponse{Part: part}, nil
}

func (s *InventoryService) ListParts(_ context.Context, req *inventoryV1.ListPartsRequest) (*inventoryV1.ListPartsResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "validation error: %v", err)
	}

	parts, err := s.storage.Parts(req.GetFilter())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list parts: %v", err)
	}

	return &inventoryV1.ListPartsResponse{Parts: parts}, nil
}

func matchesFilter(part *inventoryV1.Part, filter *inventoryV1.PartsFilter) bool {
	if len(filter.GetUuids()) > 0 && !slices.Contains(filter.GetUuids(), part.GetUuid()) {
		return false
	}

	// Фильтрация по имени
	if len(filter.GetNames()) > 0 && !slices.Contains(filter.GetNames(), part.GetName()) {
		return false
	}

	// Фильтрация по категории
	if len(filter.GetCategories()) > 0 && !slices.Contains(filter.GetCategories(), part.GetCategory()) {
		return false
	}

	// Фильтрация по странам
	if len(filter.GetManufacturerCountries()) > 0 && !slices.Contains(filter.GetManufacturerCountries(), part.GetManufacturer().GetCountry()) {
		return false
	}

	// Фильтрация по тегам (если хотя бы один тег совпадает)
	if len(filter.GetTags()) > 0 && !hasCommonElement(filter.GetTags(), part.GetTags()) {
		return false
	}

	return true
}

func hasCommonElement(a, b []string) bool {
	for _, v := range a {
		if slices.Contains(b, v) {
			return true
		}
	}
	return false
}

func generateParts() []*inventoryV1.Part {
	names := []string{
		"Main Engine",
		"Reserve Engine",
		"Thruster",
		"Fuel Tank",
		"Left Wing",
		"Right Wing",
		"Window A",
		"Window B",
		"Control Module",
		"Stabilizer",
	}

	descriptions := []string{
		"Primary propulsion unit",
		"Backup propulsion unit",
		"Thruster for fine adjustments",
		"Main fuel tank",
		"Left aerodynamic wing",
		"Right aerodynamic wing",
		"Front viewing window",
		"Side viewing window",
		"Flight control module",
		"Stabilization fin",
	}

	var parts []*inventoryV1.Part
	for i := 0; i < gofakeit.Number(1, 50); i++ {
		idx := gofakeit.Number(0, len(names)-1)
		parts = append(parts, &inventoryV1.Part{
			Uuid:          uuid.NewString(),
			Name:          names[idx],
			Description:   descriptions[idx],
			Price:         roundTo(gofakeit.Float64Range(100, 10_000)),
			StockQuantity: int64(gofakeit.Number(1, 100)),
			Category:      inventoryV1.Category(gofakeit.Number(1, 4)), //nolint:gosec // safe: gofakeit.Number returns 1..4
			Dimensions:    generateDimensions(),
			Manufacturer:  generateManufacturer(),
			Tags:          generateTags(),
			Metadata:      generateMetadata(),
			CreatedAt:     timestamppb.Now(),
		})
	}

	return parts
}

func generateDimensions() *inventoryV1.Dimensions {
	return &inventoryV1.Dimensions{
		Length: roundTo(gofakeit.Float64Range(1, 1000)),
		Width:  roundTo(gofakeit.Float64Range(1, 1000)),
		Height: roundTo(gofakeit.Float64Range(1, 1000)),
		Weight: roundTo(gofakeit.Float64Range(1, 1000)),
	}
}

func generateManufacturer() *inventoryV1.Manufacturer {
	return &inventoryV1.Manufacturer{
		Name:    gofakeit.Name(),
		Country: gofakeit.Country(),
		Website: gofakeit.URL(),
	}
}

func generateTags() []string {
	var tags []string
	for i := 0; i < gofakeit.Number(1, 10); i++ {
		tags = append(tags, gofakeit.EmojiTag())
	}

	return tags
}

func generateMetadata() map[string]*inventoryV1.Value {
	metadata := make(map[string]*inventoryV1.Value)

	for i := 0; i < gofakeit.Number(1, 10); i++ {
		metadata[gofakeit.Word()] = generateMetadataValue()
	}

	return metadata
}

func generateMetadataValue() *inventoryV1.Value {
	switch gofakeit.Number(0, 3) {
	case 0:
		return &inventoryV1.Value{
			Kind: &inventoryV1.Value_StringValue{
				StringValue: gofakeit.Word(),
			},
		}

	case 1:
		return &inventoryV1.Value{
			Kind: &inventoryV1.Value_Int64Value{
				Int64Value: int64(gofakeit.Number(1, 100)),
			},
		}

	case 2:
		return &inventoryV1.Value{
			Kind: &inventoryV1.Value_DoubleValue{
				DoubleValue: roundTo(gofakeit.Float64Range(1, 100)),
			},
		}

	case 3:
		return &inventoryV1.Value{
			Kind: &inventoryV1.Value_BoolValue{
				BoolValue: gofakeit.Bool(),
			},
		}

	default:
		return nil
	}
}

func roundTo(x float64) float64 {
	return math.Round(x*100) / 100
}
