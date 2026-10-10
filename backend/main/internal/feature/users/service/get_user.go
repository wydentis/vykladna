package users_service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	core_domains "github.com/wydentis/vykladna/shared/core/domains"
)

func (s *UsersService) GetUser(ctx context.Context, id uuid.UUID) (core_domains.User, error) {
	user, err := s.usersRepository.GetUserByID(ctx, id)
	if err != nil {
		return core_domains.User{}, fmt.Errorf("get user: %w", err)
	}

	return user, nil
}
