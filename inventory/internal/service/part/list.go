package part

import (
	"context"

	"go.uber.org/zap"

	"github.com/Alexey-step/rocket-factory/inventory/internal/model"
	"github.com/Alexey-step/rocket-factory/platform/pkg/logger"
)

func (r *service) ListParts(ctx context.Context, filter model.PartsFilter) ([]model.Part, error) {
	parts, err := r.inventoryRepository.ListParts(ctx, filter)
	if err != nil {
		logger.Error(ctx, "failed to get parts",
			zap.Any("filter", filter),
			zap.Error(err),
		)
		return nil, err
	}

	return parts, nil
}
