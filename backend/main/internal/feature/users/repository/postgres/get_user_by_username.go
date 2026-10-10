package users_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	core_postgres_pool "github.com/wydentis/vykladna/shared/core/db/postgres/pool"
	core_domains "github.com/wydentis/vykladna/shared/core/domains"
	core_errors "github.com/wydentis/vykladna/shared/utils/errors"
)

func (r *UsersRepository) GetUserByUsername(ctx context.Context, username string) (core_domains.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		SELECT id, version, role, username, name, surname, password_hash
		FROM users.accounts
		WHERE username=$1
	`

	row := r.pool.QueryRow(ctx, query, username)

	var userModel UserModel
	if err := row.Scan(
		&userModel.ID,
		&userModel.Version,
		&userModel.Role,
		&userModel.Username,
		&userModel.Name,
		&userModel.Surname,
		&userModel.PasswordHash,
	); err != nil {
		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return core_domains.User{}, fmt.Errorf("user with username=%s: %w", username, core_errors.ErrNotFound)
		}
		return core_domains.User{}, fmt.Errorf("scan error: %w", err)
	}

	userDomain := userDomainFromModel(userModel)
	return userDomain, nil
}
