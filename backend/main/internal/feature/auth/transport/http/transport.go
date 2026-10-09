package auth_transport_http

import (
	"context"
	"fmt"
	"net/http"

	core_auth "github.com/wydentis/vykladna/shared/core/auth"
	core_http_server "github.com/wydentis/vykladna/shared/core/transport_http/server"
)

var (
	refreshCookieName = "refresh_token"
	refreshCookiePath = fmt.Sprintf("/api/%s/auth", core_http_server.APIVersion1)
)

type AuthTransportHTTP struct {
	authService AuthService
}

type AuthService interface {
	GetAccessToken(ctx context.Context, refreshTok string) (string, error)
	GetAuthTokens(ctx context.Context, username, password string) (core_auth.Tokens, error)
}

func NewAuthTransportHTTP(authService AuthService) *AuthTransportHTTP {
	return &AuthTransportHTTP{
		authService: authService,
	}
}

func (h *AuthTransportHTTP) Routes() []core_http_server.Route {
	return []core_http_server.Route{
		{
			Method:  http.MethodPost,
			Path:    "/auth",
			Handler: h.GetAuthTokens,
		},
		{
			Method:  http.MethodPost,
			Path:    "/auth/refresh",
			Handler: h.GetAccessToken,
		},
		{
			Method:  http.MethodPost,
			Path:    "/auth/logout",
			Handler: h.Logout,
		},
	}
}
