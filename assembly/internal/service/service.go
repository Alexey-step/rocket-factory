package service

import (
	"context"

	"github.com/Alexey-step/rocket-factory/assembly/internal/model"
)

type OrderConsumerService interface {
	RunConsumer(ctx context.Context) error
}

type OrderProducerService interface {
	ProduceShipAssembled(ctx context.Context, event model.ShipAssembled) error
}

type AssemblyService interface{}

type service struct {
	assemblyProducerService OrderProducerService
}

func NewService(
	assemblyProducerService OrderProducerService,
) *service {
	return &service{
		assemblyProducerService: assemblyProducerService,
	}
}
