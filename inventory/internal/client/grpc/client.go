package grpc

import (
	"context"

	authV1 "github.com/Alexey-step/rocket-factory/shared/pkg/proto/auth/v1"
)

type IamClient interface {
	Login(ctx context.Context, req *authV1.LoginRequest) (*authV1.LoginResponse, error)
	Whoami(ctx context.Context, req *authV1.WhoamiRequest) (*authV1.WhoamiResponse, error)
}
