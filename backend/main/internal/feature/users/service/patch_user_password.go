package users_service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	utils_crypt "github.com/wydentis/vykladna/shared/utils/crypt"
)

func (s *UsersService) PatchUserPassword(ctx context.Context, id uuid.UUID, password string) error {
	user, err := s.usersRepository.GetUserByID(ctx, id)
	if err != nil {
		return fmt.Errorf("get user: %w", err)
	}

	hash, err := utils_crypt.HashPassword(password)
	if err != nil {
		return fmt.Errorf("hash password: %w", err)
	}

	user.PasswordHash = hash

	_, err = s.usersRepository.PatchUser(ctx, user)
	if err != nil {
		return fmt.Errorf("patch user: %w", err)
	}

	return nil
}
