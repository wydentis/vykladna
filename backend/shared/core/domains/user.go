package core_domains

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	utils_validation "github.com/wydentis/vykladna/shared/utils/validation"
)

type UserRole string

const (
	UserRoleAdmin   UserRole = "admin"
	UserRoleTeacher UserRole = "teacher"
)

var (
	ErrInvalidRole = errors.New("invalid role")
)

func (r UserRole) IsValid() bool {
	switch r {
	case UserRoleAdmin, UserRoleTeacher:
		return true
	}

	return false
}

func ParseRole(s string) (UserRole, error) {
	r := UserRole(s)
	if !r.IsValid() {
		return "", fmt.Errorf("%w: %q", ErrInvalidRole, r)
	}

	return r, nil
}

type User struct {
	ID      uuid.UUID
	Version int64

	Role         UserRole
	Username     string
	Name         string
	Surname      string
	PasswordHash string
}

func NewUser(
	ID uuid.UUID,
	Version int64,
	Role UserRole,
	Username string,
	Name string,
	Surname string,
	PasswordHash string,
) User {
	return User{
		ID:           ID,
		Version:      Version,
		Role:         Role,
		Username:     Username,
		Name:         Name,
		Surname:      Surname,
		PasswordHash: PasswordHash,
	}
}

func NewUserUninitialized(
	Username string,
	Name string,
	Surname string,
) User {
	return NewUser(
		UninitializedID,
		UninitializedVersion,
		UserRoleTeacher,
		Username,
		Name,
		Surname,
		UninitializedPassword,
	)
}

func (u *User) Validate() error {
	var subErr error

	if !u.Role.IsValid() {
		subErr = fmt.Errorf("'role' must be in {'teacher', 'admin'}")
	}
	if err := utils_validation.ValidateUsername(u.Username); err != nil {
		subErr = err
	}
	if err := utils_validation.ValidateName(u.Name); err != nil {
		subErr = err
	}
	if err := utils_validation.ValidateSurname(u.Surname); err != nil {
		subErr = err
	}

	if subErr != nil {
		return fmt.Errorf("'user' validation failed: %w", subErr)
	}

	return nil
}
