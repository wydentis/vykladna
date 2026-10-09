package auth_service

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	core_errors "github.com/wydentis/vykladna/shared/utils/errors"
)

func (s *AuthService) GetAccessToken(ctx context.Context, refreshTok string) (string, error) {
	c, err := s.authManager.ParseRefresh(refreshTok)
	if err != nil {
		return "", fmt.Errorf("parse refresh token: %v: %w", err, core_errors.ErrUnauthorized)
	}

	userID, err := uuid.Parse(c.Subject)
	if err != nil {
		return "", fmt.Errorf("parse uuid from token: %v: %w", err, core_errors.ErrUnauthorized)
	}

	user, err := s.usersReader.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, core_errors.ErrNotFound) {
			return "", fmt.Errorf("user %s not found: %w", userID, core_errors.ErrUnauthorized)
		}
		return "", fmt.Errorf("get user by id: %w", err)
	}

	return s.authManager.IssueAccess(user.ID, string(user.Role))
}
