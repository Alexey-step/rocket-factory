package user

import (
	"context"

	"golang.org/x/crypto/bcrypt"

	"github.com/Alexey-step/rocket-factory/iam/internal/model"
)

func (s *service) Register(ctx context.Context, userInfo model.UserInfo, password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	userUUID, err := s.userRepository.Create(ctx, userInfo, hashedPassword)
	if err != nil {
		return "", err
	}

	return userUUID, nil
}
