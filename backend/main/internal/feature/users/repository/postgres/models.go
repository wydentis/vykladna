package users_postgres_repository

import (
	"github.com/google/uuid"
	core_domains "github.com/wydentis/vykladna/shared/core/domains"
)

type UserModel struct {
	ID      uuid.UUID
	Version int64

	Role         core_domains.UserRole
	Username     string
	Name         string
	Surname      string
	PasswordHash string
}

func userDomainFromModel(model UserModel) core_domains.User {
	return core_domains.NewUser(
		model.ID,
		model.Version,
		model.Role,
		model.Username,
		model.Name,
		model.Surname,
		model.PasswordHash,
	)
}
