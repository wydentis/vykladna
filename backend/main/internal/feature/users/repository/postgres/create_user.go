package users_postgres_repository

import (
	"context"
	"fmt"

	core_domains "github.com/wydentis/vykladna/shared/core/domains"
)

func (r *UsersRepository) CreateUser(ctx context.Context, user core_domains.User) (core_domains.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		INSERT INTO users.accounts (role, username, name, surname, password_hash)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, version, role, username, name, surname, password_hash
	`

	row := r.pool.QueryRow(ctx, query, user.Role, user.Username, user.Name, user.Surname, user.PasswordHash)

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
		return core_domains.User{}, fmt.Errorf("scan error: %w", err)
	}

	return userDomainFromModel(userModel), nil
}
