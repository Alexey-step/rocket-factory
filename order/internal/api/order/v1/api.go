package v1

import (
	"github.com/Alexey-step/rocket-factory/order/internal/service"
)

type api struct {
	service service.OrderService
}

func NewAPI(service service.OrderService) *api {
	return &api{
		service: service,
	}
}
