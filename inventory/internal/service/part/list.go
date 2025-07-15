package part

import (
	"context"

	"github.com/Alexey-step/rocket-factory/inventory/internal/model"
)

func (r *service) ListParts(ctx context.Context, filter model.PartsFilter) ([]model.Part, error) {
	parts, err := r.inventoryRepository.ListParts(ctx, filter)
	if err != nil {
		return nil, err
	}

	return parts, nil
}
