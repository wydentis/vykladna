package users_transport_http

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	core_auth "github.com/wydentis/vykladna/shared/core/auth"
	core_domains "github.com/wydentis/vykladna/shared/core/domains"
	core_http_middleware "github.com/wydentis/vykladna/shared/core/transport_http/middleware"
	core_http_server "github.com/wydentis/vykladna/shared/core/transport_http/server"
)

type UsersHTTPTransport struct {
	usersService UsersService
	authManager  *core_auth.Manager
}

type UsersService interface {
	GetUser(ctx context.Context, id uuid.UUID) (core_domains.User, error)
	CreateUser(ctx context.Context, user core_domains.User, password string) (core_domains.User, error)
	PatchUser(ctx context.Context, id uuid.UUID, userPatch core_domains.UserPatch) (core_domains.User, error)
	PatchUserPassword(ctx context.Context, id uuid.UUID, password string) error
}

func NewUsersHTTPTransport(userService UsersService, authManager *core_auth.Manager) *UsersHTTPTransport {
	return &UsersHTTPTransport{
		usersService: userService,
		authManager:  authManager,
	}
}

func (h *UsersHTTPTransport) Routes() []core_http_server.Route {
	return []core_http_server.Route{
		{
			Method:      http.MethodGet,
			Path:        "/users/me",
			Handler:     h.GetMe,
			Middlewares: []core_http_middleware.Middleware{core_http_middleware.Auth(*h.authManager)},
		},
		{
			Method:  http.MethodPost,
			Path:    "/users",
			Handler: h.CreateUser,
		},
		{
			Method:      http.MethodPatch,
			Path:        "/user",
			Handler:     h.PatchUser,
			Middlewares: []core_http_middleware.Middleware{core_http_middleware.Auth(*h.authManager)},
		},
		{
			Method:      http.MethodPatch,
			Path:        "/user/password",
			Handler:     h.PatchUserPassword,
			Middlewares: []core_http_middleware.Middleware{core_http_middleware.Auth(*h.authManager)},
		},
	}
}
