package core_auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const (
	audAccess  = "vykladna-access"
	audRefresh = "vykladna-refresh"
)

type Manager struct {
	secret     []byte
	issuer     string
	AccessTTL  time.Duration
	RefreshTTL time.Duration
}

func NewManager(config Config) *Manager {
	return &Manager{
		secret:     []byte(config.Secret),
		issuer:     config.Issuer,
		AccessTTL:  config.AccessTTL,
		RefreshTTL: config.RefreshTTL,
	}
}

func (m *Manager) sign(userID uuid.UUID, role, aud string, ttl time.Duration) (string, error) {
	now := time.Now()
	c := Claims{
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject: userID.String(),
			Issuer:  m.issuer,
			Audience: jwt.ClaimStrings{
				aud,
			},
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
			ID:        uuid.NewString(),
		},
	}

	return jwt.NewWithClaims(jwt.SigningMethodHS256, c).SignedString(m.secret)
}

func (m *Manager) parse(tok, aud string) (*Claims, error) {
	var c Claims
	_, err := jwt.ParseWithClaims(
		tok,
		&c,
		func(t *jwt.Token) (any, error) { return m.secret, nil },
		jwt.WithValidMethods([]string{"HS256"}),
		jwt.WithIssuer(m.issuer),
		jwt.WithAudience(aud),
		jwt.WithExpirationRequired(),
		jwt.WithIssuedAt(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	return &c, nil
}

func (m *Manager) IssueAccess(userID uuid.UUID, role string) (string, error) {
	return m.sign(userID, role, audAccess, m.AccessTTL)
}

func (m *Manager) IssueRefresh(userID uuid.UUID) (string, error) {
	return m.sign(userID, "", audRefresh, m.RefreshTTL)
}

func (m *Manager) ParseAccess(tok string) (*Claims, error) {
	return m.parse(tok, audAccess)
}

func (m *Manager) ParseRefresh(tok string) (*Claims, error) {
	return m.parse(tok, audRefresh)
}
