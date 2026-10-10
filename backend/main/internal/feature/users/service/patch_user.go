package users_service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	core_domains "github.com/wydentis/vykladna/shared/core/domains"
	core_errors "github.com/wydentis/vykladna/shared/utils/errors"
)

func (s *UsersService) PatchUser(ctx context.Context, id uuid.UUID, patch core_domains.UserPatch) (core_domains.User, error) {
	user, err := s.usersRepository.GetUserByID(ctx, id)
	if err != nil {
		return core_domains.User{}, fmt.Errorf("get user: %w", err)
	}

	if err := user.ApplyPatch(patch); err != nil {
		return core_domains.User{}, fmt.Errorf("apply user patch: %v: %w", err, core_errors.ErrInvalidArgument)
	}

	user, err = s.usersRepository.PatchUser(ctx, user)
	if err != nil {
		return core_domains.User{}, fmt.Errorf("patch user: %w", err)
	}

	return user, nil
}
