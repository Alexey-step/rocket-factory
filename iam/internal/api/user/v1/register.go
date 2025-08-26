package v1

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/Alexey-step/rocket-factory/iam/internal/converter"
	"github.com/Alexey-step/rocket-factory/platform/pkg/logger"
	userV1 "github.com/Alexey-step/rocket-factory/shared/pkg/proto/user/v1"
)

func (a *api) Register(ctx context.Context, req *userV1.RegisterRequest) (*userV1.RegisterResponse, error) {
	userInfo := converter.UserInfoToModel(req.GetInfo().Info)
	userUUID, err := a.service.Register(ctx, userInfo, req.GetInfo().Password)
	if err != nil {
		logger.Error(ctx, "error while registering user",
			zap.Error(err),
		)
		return nil, status.Errorf(codes.Internal, "internal error while registering user")
	}

	return &userV1.RegisterResponse{
		UserUuid: userUUID,
	}, nil
}
