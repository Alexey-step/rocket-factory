package v1

import (
	def "github.com/Alexey-step/rocket-factory/order/internal/client/grpc"
	generatedAuthV1 "github.com/Alexey-step/rocket-factory/shared/pkg/proto/auth/v1"
	generatedUserV1 "github.com/Alexey-step/rocket-factory/shared/pkg/proto/user/v1"
)

var _ def.IamClient = (*client)(nil)

type client struct {
	generatedAuthClient generatedAuthV1.AuthServiceClient
	generatedUserV1     generatedUserV1.UserServiceClient
}

func NewClient(
	generatedAuthClient generatedAuthV1.AuthServiceClient,
	generatedUserV1 generatedUserV1.UserServiceClient,
) *client {
	return &client{
		generatedAuthClient: generatedAuthClient,
		generatedUserV1:     generatedUserV1,
	}
}
