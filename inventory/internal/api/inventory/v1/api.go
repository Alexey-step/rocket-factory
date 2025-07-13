package v1

import (
	"github.com/Alexey-step/rocket-factory/inventory/internal/service"
	inventory_v1 "github.com/Alexey-step/rocket-factory/shared/pkg/proto/inventory/v1"
)

type api struct {
	inventory_v1.UnimplementedInventoryServiceServer

	service service.InventoryService
}

func NewAPI(service service.InventoryService) *api {
	return &api{
		service: service,
	}
}
