package v1

import (
	"context"
	"errors"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/Alexey-step/rocket-factory/iam/internal/converter"
	"github.com/Alexey-step/rocket-factory/iam/internal/model"
	"github.com/Alexey-step/rocket-factory/platform/pkg/logger"
	userV1 "github.com/Alexey-step/rocket-factory/shared/pkg/proto/user/v1"
)

func (a *api) GetUser(ctx context.Context, req *userV1.GetUserRequest) (*userV1.GetUserResponse, error) {
	user, err := a.service.GetUser(ctx, req.GetUserUuid())
	if err != nil {
		if errors.Is(err, model.ErrUserNotFound) {
			logger.Error(ctx, "error while getting user",
				zap.String("user_uuid", req.GetUserUuid()),
				zap.Error(err),
			)
			return nil, status.Errorf(codes.NotFound, "user with UUID %s not found", req.GetUserUuid())
		}
		logger.Error(ctx, "error while getting user",
			zap.String("user_uuid", req.GetUserUuid()),
			zap.Error(err),
		)
		return nil, status.Errorf(codes.Internal, "internal error while getting user")
	}

	return &userV1.GetUserResponse{
		User: converter.UserToProto(user),
	}, nil
}
