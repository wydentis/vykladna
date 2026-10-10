package users_service

import (
	"context"

	"github.com/google/uuid"
	core_domains "github.com/wydentis/vykladna/shared/core/domains"
)

type UsersService struct {
	usersRepository UsersRepository
}

type UsersRepository interface {
	GetUserByID(ctx context.Context, id uuid.UUID) (core_domains.User, error)
	GetUserByUsername(ctx context.Context, username string) (core_domains.User, error)
	CreateUser(ctx context.Context, user core_domains.User) (core_domains.User, error)
}

func NewUsersService(usersRepository UsersRepository) *UsersService {
	return &UsersService{
		usersRepository: usersRepository,
	}
}
