package part

import (
	"context"

	"github.com/Alexey-step/rocket-factory/inventory/internal/model"
	repoConverter "github.com/Alexey-step/rocket-factory/inventory/internal/repository/converter"
)

func (r *repository) GetPart(_ context.Context, orderUUID string) (model.Part, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	repoPart, ok := r.data[orderUUID]
	if !ok {
		return model.Part{}, model.ErrPartNotFound
	}

	return repoConverter.PartToModel(repoPart), nil
}
