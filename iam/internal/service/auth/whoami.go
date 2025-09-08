package auth

import (
	"context"

	"github.com/Alexey-step/rocket-factory/iam/internal/model"
)

func (s *service) Whoami(ctx context.Context, sessionUUID string) (model.Session, model.User, error) {
	session, user, err := s.sessionRepository.Get(ctx, sessionUUID)
	if err != nil {
		return model.Session{}, model.User{}, err
	}

	return session, user, nil
}
