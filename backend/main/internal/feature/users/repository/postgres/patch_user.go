package users_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	core_postgres_pool "github.com/wydentis/vykladna/shared/core/db/postgres/pool"
	core_domains "github.com/wydentis/vykladna/shared/core/domains"
	core_errors "github.com/wydentis/vykladna/shared/utils/errors"
)

func (r *UsersRepository) PatchUser(ctx context.Context, user core_domains.User) (core_domains.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	UPDATE users.accounts
	SET name = $1, surname = $2, password_hash = $3, version = version + 1
	WHERE id = $4 AND version = $5
	RETURNING id, version, role, username, name, surname, password_hash
	`

	row := r.pool.QueryRow(
		ctx, query,
		user.Name,
		user.Surname,
		user.PasswordHash,
		user.ID,
		user.Version,
	)

	var model UserModel
	if err := row.Scan(
		&model.ID,
		&model.Version,
		&model.Role,
		&model.Username,
		&model.Name,
		&model.Surname,
		&model.PasswordHash,
	); err != nil {
		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return core_domains.User{}, fmt.Errorf("user with id=%s concurrently accessed: %w", user.ID, core_errors.ErrConflict)
		}

		return core_domains.User{}, fmt.Errorf("scan error: %w", err)
	}

	domain := userDomainFromModel(model)
	return domain, nil
}
