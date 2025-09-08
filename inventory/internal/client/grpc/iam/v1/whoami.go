package v1

import (
	"context"

	generatedAuthV1 "github.com/Alexey-step/rocket-factory/shared/pkg/proto/auth/v1"
)

func (c *client) Whoami(ctx context.Context, req *generatedAuthV1.WhoamiRequest) (*generatedAuthV1.WhoamiResponse, error) {
	res, err := c.generatedAuthClient.Whoami(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
