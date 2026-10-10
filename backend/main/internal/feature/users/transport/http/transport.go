package users_transport_http

import (
	"context"
	"net/http"

	core_domains "github.com/wydentis/vykladna/shared/core/domains"
	core_http_server "github.com/wydentis/vykladna/shared/core/transport_http/server"
)

type UsersHTTPTransport struct {
	usersService UsersService
}

type UsersService interface {
	CreateUser(ctx context.Context, user core_domains.User, password string) (core_domains.User, error)
}

func NewUsersHTTPTransport(userService UsersService) *UsersHTTPTransport {
	return &UsersHTTPTransport{
		usersService: userService,
	}
}

func (h *UsersHTTPTransport) Routes() []core_http_server.Route {
	return []core_http_server.Route{
		{
			Method:  http.MethodPost,
			Path:    "/users",
			Handler: h.CreateUser,
		},
	}
}
