package v1

import (
	"context"

	clientConverter "github.com/Alexey-step/rocket-factory/order/internal/client/converter"
	"github.com/Alexey-step/rocket-factory/order/internal/model"
	grpcAuth "github.com/Alexey-step/rocket-factory/platform/pkg/middleware/grpc"
	generatedInventoryV1 "github.com/Alexey-step/rocket-factory/shared/pkg/proto/inventory/v1"
)

func (c *client) ListParts(ctx context.Context, filter model.PartsFilter) (parts []model.Part, error error) {
	// пробрасываем данные из outcoming context в outgoing gRPC context
	ctx = grpcAuth.ForwardSessionUUIDToGRPC(ctx)

	partsFilter := &generatedInventoryV1.ListPartsRequest{
		Filter: clientConverter.PartsFilterToProto(filter),
	}

	partsList, err := c.generatedClient.ListParts(ctx, partsFilter)
	if err != nil {
		return nil, err
	}

	return clientConverter.PartListToModel(partsList.Parts), nil
}
