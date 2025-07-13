package part

import (
	"context"

	"github.com/Alexey-step/rocket-factory/inventory/internal/model"
)

func (r *service) GetPart(ctx context.Context, orderUUID string) (model.Part, error) {
	part, err := r.inventoryRepository.GetPart(ctx, orderUUID)
	if err != nil {
		return model.Part{}, err
	}

	return part, nil
}
