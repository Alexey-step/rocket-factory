package part

import (
	def "github.com/Alexey-step/rocket-factory/inventory/internal/service"
)

var _ def.InventoryService = (*service)(nil)

type service struct {
	inventoryRepository def.InventoryRepository
}

func NewService(
	inventoryRepository def.InventoryRepository,
) *service {
	return &service{
		inventoryRepository: inventoryRepository,
	}
}
