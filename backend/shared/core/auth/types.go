package core_auth

import (
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	Role string `json:"role"`
	jwt.RegisteredClaims
}

type Tokens struct {
	Access  string
	Refresh string
}

func NewTokens(access, refresh string) Tokens {
	return Tokens{
		Access:  access,
		Refresh: refresh,
	}
}
