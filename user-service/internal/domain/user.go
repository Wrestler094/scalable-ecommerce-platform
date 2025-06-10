package domain

import (
	"context"

	"pkg/authenticator"
)

// User представляет публичного пользователя, безопасного для отдачи в API или кеширования.
type User struct {
	ID    int64
	Email string
	Role  authenticator.Role
}

// HashedPassword — отдельный тип для хранения хэша пароля, чтобы избежать случайного использования.
type HashedPassword string

// UserWithPassword представляет внутреннего пользователя, содержащего чувствительные данные (например, хэш пароля).
type UserWithPassword struct {
	User
	Password HashedPassword
}

type PasswordHasher interface {
	Hash(password string) (string, error)
	Compare(hashed, plain string) bool
}

type TokenManager interface {
	GenerateAccessToken(userID int64, role string) (string, error)
	GenerateRefreshToken() (string, error)
}

type UserRepository interface {
	CreateUser(ctx context.Context, user UserWithPassword) (int64, error)
	GetUserByEmail(ctx context.Context, email string) (*UserWithPassword, error)
	GetUserByID(ctx context.Context, id int64) (*User, error)
}

type RefreshTokenRepository interface {
	Store(ctx context.Context, userID int64, token string) error
	GetUserID(ctx context.Context, token string) (int64, error)
	Delete(ctx context.Context, token string) error
	Replace(ctx context.Context, oldToken string, newToken string, userID int64) error
}
