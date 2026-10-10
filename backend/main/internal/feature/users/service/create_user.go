package users_service

import (
	"context"
	"errors"
	"fmt"

	core_domains "github.com/wydentis/vykladna/shared/core/domains"
	utils_crypt "github.com/wydentis/vykladna/shared/utils/crypt"
	core_errors "github.com/wydentis/vykladna/shared/utils/errors"
	utils_validation "github.com/wydentis/vykladna/shared/utils/validation"
)

func (s *UsersService) CreateUser(ctx context.Context, user core_domains.User, password string) (core_domains.User, error) {
	_, err := s.usersRepository.GetUserByUsername(ctx, user.Username)
	switch {
	case err == nil:
		return core_domains.User{}, fmt.Errorf("user with username=%q is already exists: %w", user.Username, core_errors.ErrConflict)
	case !errors.Is(err, core_errors.ErrNotFound):
		return core_domains.User{}, fmt.Errorf("get user: %w", err)
	}

	if err := utils_validation.ValidatePassword(password); err != nil {
		return core_domains.User{}, fmt.Errorf("validate 'password': %w", err)
	}

	pswHash, err := utils_crypt.HashPassword(password)
	if err != nil {
		return core_domains.User{}, fmt.Errorf("get hash: %w", err)
	}

	user.PasswordHash = pswHash

	if err := user.Validate(); err != nil {
		return core_domains.User{}, fmt.Errorf("validate user: %v: %w", err, core_errors.ErrInvalidArgument)
	}

	user, err = s.usersRepository.CreateUser(ctx, user)
	if err != nil {
		return core_domains.User{}, fmt.Errorf("create user: %w", err)
	}

	return user, err
}
