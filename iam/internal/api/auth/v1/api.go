package v1

import (
	authv3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"

	"github.com/Alexey-step/rocket-factory/iam/internal/service"
	authV1 "github.com/Alexey-step/rocket-factory/shared/pkg/proto/auth/v1"
)

type api struct {
	authV1.UnimplementedAuthServiceServer
	authv3.UnimplementedAuthorizationServer

	service service.AuthService
}

func NewAPI(service service.AuthService) *api {
	return &api{
		service: service,
	}
}
