package order

import (
	"sync"

	def "github.com/Alexey-step/rocket-factory/order/internal/repository"
	repoModel "github.com/Alexey-step/rocket-factory/order/internal/repository/model"
)

var _ def.OrderRepository = (*repository)(nil)

type repository struct {
	mu     sync.RWMutex
	orders map[string]repoModel.OrderData
}

func NewOrderRepository() *repository {
	return &repository{
		orders: make(map[string]repoModel.OrderData),
	}
}
