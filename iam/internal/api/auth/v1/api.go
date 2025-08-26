package v1

import (
	"github.com/Alexey-step/rocket-factory/iam/internal/service"
	authV1 "github.com/Alexey-step/rocket-factory/shared/pkg/proto/auth/v1"
)

type api struct {
	authV1.UnimplementedAuthServiceServer

	service service.AuthService
}

func NewAPI(service service.AuthService) *api {
	return &api{
		service: service,
	}
}
