package part

import (
	"sync"

	def "github.com/Alexey-step/rocket-factory/inventory/internal/repository"
	repoModel "github.com/Alexey-step/rocket-factory/inventory/internal/repository/model"
)

var _ def.InventoryRepository = (*repository)(nil)

type repository struct {
	mu   sync.RWMutex
	data map[string]repoModel.Part
}

func NewRepository() *repository {
	s := &repository{
		data: make(map[string]repoModel.Part),
	}

	s.InitParts()

	return s
}
