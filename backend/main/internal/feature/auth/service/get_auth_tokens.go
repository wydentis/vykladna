package auth_service

import (
	"context"
	"errors"
	"fmt"

	core_auth "github.com/wydentis/vykladna/shared/core/auth"
	utils_crypt "github.com/wydentis/vykladna/shared/utils/crypt"
	core_errors "github.com/wydentis/vykladna/shared/utils/errors"
)

func (s *AuthService) GetAuthTokens(ctx context.Context, username, password string) (core_auth.Tokens, error) {
	user, err := s.usersReader.GetUserByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, core_errors.ErrNotFound) {
			return core_auth.Tokens{}, fmt.Errorf("get user: %v: %w", err, core_errors.ErrUnauthorized)
		}
		return core_auth.Tokens{}, fmt.Errorf("get user: %w", err)
	}

	if err := utils_crypt.CheckPasswordHash(password, user.PasswordHash); err != nil {
		return core_auth.Tokens{}, fmt.Errorf("hash password: %w", err)
	}

	access, err := s.authManager.IssueAccess(user.ID, string(user.Role))
	if err != nil {
		return core_auth.Tokens{}, fmt.Errorf("generate access token: %w", err)
	}

	refresh, err := s.authManager.IssueRefresh(user.ID)
	if err != nil {
		return core_auth.Tokens{}, fmt.Errorf("generate refresh token: %w", err)
	}

	tokens := core_auth.NewTokens(access, refresh)
	return tokens, nil
}
