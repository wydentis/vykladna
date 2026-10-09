package auth_transport_http

import (
	"net/http"
)

type TokenDTO struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}

func tokenDTOFromAccess(access string) TokenDTO {
	return TokenDTO{
		AccessToken: access,
		TokenType:   "Bearer",
	}
}

func newRefreshCookie(value string, maxAge int) *http.Cookie {
	return &http.Cookie{
		Name:     refreshCookieName,
		Value:    value,
		Path:     refreshCookiePath,
		MaxAge:   maxAge,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	}
}
