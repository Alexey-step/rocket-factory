package v1

import (
	"context"

	"github.com/Alexey-step/rocket-factory/iam/internal/converter"
	authV1 "github.com/Alexey-step/rocket-factory/shared/pkg/proto/auth/v1"
)

func (a *api) Whoami(ctx context.Context, req *authV1.WhoamiRequest) (*authV1.WhoamiResponse, error) {
	session, user, err := a.service.Whoami(ctx, req.SessionUuid)
	if err != nil {
		return nil, err
	}

	return &authV1.WhoamiResponse{
		Session: converter.SessionToProto(session),
		User:    converter.UserToProto(user),
	}, nil
}
