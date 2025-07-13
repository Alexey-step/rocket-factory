package v1

import (
	"context"

	"github.com/Alexey-step/rocket-factory/order/internal/client/converter"
	"github.com/Alexey-step/rocket-factory/order/internal/model"
	generatedInventoryV1 "github.com/Alexey-step/rocket-factory/shared/pkg/proto/inventory/v1"
)

func (c *client) ListParts(ctx context.Context, filter model.PartsFilter) (parts []model.Part, error error) {
	partsList, err := c.generatedClient.ListParts(ctx, &generatedInventoryV1.ListPartsRequest{
		Filter: converter.PartsFilterToProto(filter),
	})
	if err != nil {
		return nil, err
	}

	return converter.PartListToModel(partsList.Parts), nil
}
