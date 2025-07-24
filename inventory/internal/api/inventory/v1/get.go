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

func (a *api) GetPart(ctx context.Context, req *inventoryV1.GetPartRequest) (*inventoryV1.GetPartResponse, error) {
	part, err := a.service.GetPart(ctx, req.Uuid)
	if err != nil {
		if errors.Is(err, model.ErrPartNotFound) {
			return nil, status.Errorf(codes.NotFound, "part with UUID %s not found", req.GetUuid())
		}
		if errors.Is(err, model.ErrPartsInternalError) {
			return nil, status.Errorf(codes.Internal, "internal error while getting part with UUID %s", req.GetUuid())
		}
		return nil, err
	}

	return &inventoryV1.GetPartResponse{
		Part: converter.PartToProto(part),
	}, nil
}
