package v1

import (
	"context"
	"errors"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/Alexey-step/rocket-factory/iam/internal/model"
	"github.com/Alexey-step/rocket-factory/platform/pkg/logger"
	authV1 "github.com/Alexey-step/rocket-factory/shared/pkg/proto/auth/v1"
)

func (a *api) Login(ctx context.Context, req *authV1.LoginRequest) (*authV1.LoginResponse, error) {
	sessionUUID, err := a.service.Login(ctx, req.GetLogin(), req.GetPassword())
	if err != nil {
		if errors.Is(err, model.ErrSessionBadRequest) {
			logger.Error(ctx, "bad request during login", zap.Error(err))
			return nil, status.Errorf(codes.InvalidArgument, "bad request during login")
		}
		logger.Error(ctx, "failed to login", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "internal error during login")
	}

	return &authV1.LoginResponse{
		SessionUuid: sessionUUID,
	}, nil
}
