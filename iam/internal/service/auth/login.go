package auth

import (
	"context"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/Alexey-step/rocket-factory/iam/internal/model"
)

func (s *service) Login(ctx context.Context, login, password string) (string, error) {
	user, err := s.userRepository.GetUserByLogin(ctx, login, password)
	if err != nil {
		return "", model.ErrInvalidCredentials
	}

	err = bcrypt.CompareHashAndPassword(user.Password, []byte(password))
	if err != nil {
		return "", err
	}

	session := model.Session{
		UUID:      uuid.NewString(),
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(s.cacheTTL),
	}

	err = s.sessionRepository.Create(ctx, session, user, s.cacheTTL)
	if err != nil {
		return "", err
	}

	err = s.sessionRepository.AddSessionToUserSet(ctx, user.UUID, session.UUID)
	if err != nil {
		return "", err
	}

	return session.UUID, nil
}
