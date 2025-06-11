package jwt

import (
	"fmt"
	"time"

	"user-service/internal/domain"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var _ domain.TokenManager = (*Manager)(nil)

// accessTokenClaims describes the structure of JWT used between services.
// Must contain `sub` (user ID) and `role`.
type accessTokenClaims struct {
	Role string `json:"role"`
	jwt.RegisteredClaims
}

type Manager struct {
	accessSecret  string
	accessTTL     time.Duration
	signingMethod jwt.SigningMethod
}

func NewManager(accessSecret string, accessTTL time.Duration) *Manager {
	return &Manager{
		accessSecret:  accessSecret,
		accessTTL:     accessTTL,
		signingMethod: jwt.SigningMethodHS256,
	}
}

func (m *Manager) GenerateAccessToken(userID int64, role string) (string, error) {
	const op = "jwt.Manager.GenerateAccessToken"

	now := time.Now()

	claims := accessTokenClaims{
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   fmt.Sprint(userID),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(m.accessTTL)),
		},
	}

	token := jwt.NewWithClaims(m.signingMethod, claims)
	signed, err := token.SignedString([]byte(m.accessSecret))
	if err != nil {
		return "", fmt.Errorf("%s: failed to sign access token: %w", op, err)
	}

	return signed, nil
}

func (m *Manager) GenerateRefreshToken() (string, error) {
	const op = "jwt.Manager.GenerateRefreshToken"

	id, err := uuid.NewRandom()
	if err != nil {
		return "", fmt.Errorf("%s: failed to generate UUID: %w", op, err)
	}

	return id.String(), nil
}
