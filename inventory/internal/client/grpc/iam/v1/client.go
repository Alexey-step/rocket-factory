package v1

import (
	inventoryClient "github.com/Alexey-step/rocket-factory/inventory/internal/client/grpc"
	generatedAuthV1 "github.com/Alexey-step/rocket-factory/shared/pkg/proto/auth/v1"
	generatedUserV1 "github.com/Alexey-step/rocket-factory/shared/pkg/proto/user/v1"
)

var _ inventoryClient.IamClient = (*client)(nil)

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
