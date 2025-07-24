package repository

import (
	"context"

	"github.com/Alexey-step/rocket-factory/inventory/internal/model"
)

type InventoryRepository interface {
	ListParts(ctx context.Context, filter model.PartsFilter) ([]model.Part, error)
	GetPart(ctx context.Context, orderUUID string) (model.Part, error)
	InitParts()
}
