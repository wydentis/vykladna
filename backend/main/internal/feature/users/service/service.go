package users_service

import (
	"context"

	core_domains "github.com/wydentis/vykladna/shared/core/domains"
)

type UsersService struct {
	usersRepository UsersRepository
}

type UsersRepository interface {
	CreateUser(ctx context.Context, user core_domains.User) (core_domains.User, error)
}

func NewUsersService(usersRepository UsersRepository) *UsersService {
	return &UsersService{
		usersRepository: usersRepository,
	}
}
