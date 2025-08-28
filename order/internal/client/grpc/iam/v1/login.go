package v1

import (
	"context"

	generatedAuthV1 "github.com/Alexey-step/rocket-factory/shared/pkg/proto/auth/v1"
)

func (c *client) Login(ctx context.Context, req *generatedAuthV1.LoginRequest) (*generatedAuthV1.LoginResponse, error) {
	res, err := c.generatedAuthClient.Login(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
