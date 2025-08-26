package user

import (
	iamService "github.com/Alexey-step/rocket-factory/iam/internal/service"
)

var _ iamService.UserService = (*service)(nil)

type service struct {
	userRepository iamService.UserRepository
}

func NewService(
	userRepository iamService.UserRepository,
) *service {
	return &service{
		userRepository: userRepository,
	}
}
