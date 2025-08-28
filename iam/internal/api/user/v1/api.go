package v1

import (
	"github.com/Alexey-step/rocket-factory/iam/internal/service"
	userV1 "github.com/Alexey-step/rocket-factory/shared/pkg/proto/user/v1"
)

type api struct {
	userV1.UnimplementedUserServiceServer

	service service.UserService
}

func NewAPI(service service.UserService) *api {
	return &api{
		service: service,
	}
}
