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
	return token.SignedString([]byte(m.accessSecret))
}

func (m *Manager) GenerateRefreshToken() (string, error) {
	return uuid.NewString(), nil
}
