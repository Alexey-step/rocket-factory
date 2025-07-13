package v1

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/Alexey-step/rocket-factory/inventory/internal/converter"
	"github.com/Alexey-step/rocket-factory/inventory/internal/model"
	inventory_v1 "github.com/Alexey-step/rocket-factory/shared/pkg/proto/inventory/v1"
)

func (a *api) ListParts(ctx context.Context, req *inventory_v1.ListPartsRequest) (*inventory_v1.ListPartsResponse, error) {
	parts, err := a.service.ListParts(ctx, converter.PartsFilterToModel(req.Filter))
	if err != nil {
		if errors.Is(err, model.ErrPartsInternalError) {
			return nil, status.Errorf(codes.Internal, "inventory service error: %v", err)
		}
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return nil, status.Errorf(codes.Unavailable, "inventory service timeout")
		}
		return nil, err
	}

	return &inventory_v1.ListPartsResponse{
		Parts: converter.PartsToProto(parts),
	}, nil
}
