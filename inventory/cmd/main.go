package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"slices"
	"sync"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"

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

	s := grpc.NewServer()

	service := &inventoryService{
		parts: make(map[string]*inventoryV1.Part),
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

type inventoryService struct {
	inventoryV1.UnimplementedInventoryServiceServer

	mu    sync.RWMutex
	parts map[string]*inventoryV1.Part
}

func (s *inventoryService) GetPart(_ context.Context, req *inventoryV1.GetPartRequest) (*inventoryV1.GetPartResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	part, ok := s.parts[req.GetUuid()]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "part with UUID %s not found", req.GetUuid())
	}

	return &inventoryV1.GetPartResponse{
		Part: part,
	}, nil
}

func (s *inventoryService) ListParts(_ context.Context, req *inventoryV1.ListPartsRequest) (*inventoryV1.ListPartsResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []*inventoryV1.Part

	filter := req.GetFilter()

	for _, part := range s.parts {
		if !matchesFilter(part, filter) {
			continue
		}

		result = append(result, part)
	}

	return &inventoryV1.ListPartsResponse{
		Parts: result,
	}, nil
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
