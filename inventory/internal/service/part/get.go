package part

import (
	"context"

	"go.uber.org/zap"

	"github.com/Alexey-step/rocket-factory/inventory/internal/model"
	"github.com/Alexey-step/rocket-factory/platform/pkg/logger"
)

func (s *service) GetPart(ctx context.Context, orderUUID string) (model.Part, error) {
	part, err := s.inventoryRepository.GetPart(ctx, orderUUID)
	if err != nil {
		logger.Error(ctx, "failed to get part",
			zap.String("order_uuid", orderUUID),
			zap.Error(err),
		)
		return model.Part{}, err
	}

	return part, nil
}
