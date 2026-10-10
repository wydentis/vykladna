package auth_service

import (
	"context"

	"github.com/google/uuid"
	core_auth "github.com/wydentis/vykladna/shared/core/auth"
	core_domains "github.com/wydentis/vykladna/shared/core/domains"
)

type AuthService struct {
	usersReader UsersReader
	authManager *core_auth.Manager
}

type UsersReader interface {
	GetUserByID(ctx context.Context, id uuid.UUID) (core_domains.User, error)
	GetUserByUsername(ctx context.Context, username string) (core_domains.User, error)
}

func NewAuthService(usersReader UsersReader, authManager *core_auth.Manager) *AuthService {
	return &AuthService{
		usersReader: usersReader,
		authManager: authManager,
	}
}
