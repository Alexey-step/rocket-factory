package user

import (
	"context"

	"github.com/Alexey-step/rocket-factory/iam/internal/model"
)

func (s *service) GetUser(ctx context.Context, userUUID string) (model.User, error) {
	user, err := s.userRepository.GetUser(ctx, userUUID)
	if err != nil {
		return model.User{}, err
	}

	return user, nil
}
