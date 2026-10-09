package core_auth

import (
	"context"

	"github.com/google/uuid"
	core_domains "github.com/wydentis/vykladna/shared/core/domains"
)

type UserInfo struct {
	ID   uuid.UUID
	Role core_domains.UserRole
}

type UserInfoContextKey struct{}

var (
	key = UserInfoContextKey{}
)

func NewUserInfo(id uuid.UUID, role core_domains.UserRole) UserInfo {
	return UserInfo{
		ID:   id,
		Role: role,
	}
}

func (i *UserInfo) ToContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, key, i)
}

func FromContext(ctx context.Context) *UserInfo {
	userInfo, ok := ctx.Value(key).(*UserInfo)
	if !ok {
		panic("no user info in context")
	}

	return userInfo
}
