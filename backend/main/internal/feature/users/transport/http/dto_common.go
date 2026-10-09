package users_transport_http

import (
	"github.com/google/uuid"
	core_domains "github.com/wydentis/vykladna/shared/core/domains"
)

type UserDTO struct {
	ID      uuid.UUID `json:"id"`
	Version int64     `json:"version"`

	Role     core_domains.UserRole `json:"role"`
	Username string                `json:"username"`
	Name     string                `json:"name"`
	Surname  string                `json:"surname"`
}

func userDTOFromDomain(user core_domains.User) UserDTO {
	return UserDTO{
		ID:       user.ID,
		Version:  user.Version,
		Role:     user.Role,
		Username: user.Username,
		Name:     user.Name,
		Surname:  user.Surname,
	}
}
