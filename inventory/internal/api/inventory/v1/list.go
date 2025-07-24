package v1

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/Alexey-step/rocket-factory/inventory/internal/converter"
	"github.com/Alexey-step/rocket-factory/inventory/internal/model"
	inventoryV1 "github.com/Alexey-step/rocket-factory/shared/pkg/proto/inventory/v1"
)

func (a *api) ListParts(ctx context.Context, req *inventoryV1.ListPartsRequest) (*inventoryV1.ListPartsResponse, error) {
	filters := converter.PartsFilterToModel(req.GetFilter())

	parts, err := a.service.ListParts(ctx, filters)
	if err != nil {
		if errors.Is(err, model.ErrPartsNotFound) {
			return nil, status.Errorf(codes.NotFound, "inventory service error: %v", err)
		}
		return nil, status.Errorf(codes.Internal, "inventory service error: %v", err)
	}

	return &inventoryV1.ListPartsResponse{
		Parts: converter.PartsToProto(parts),
	}, nil
}
