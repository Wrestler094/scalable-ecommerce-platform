package usecase

import (
	"context"
	"errors"
	"fmt"

	"pkg/roles"
	"user-service/internal/domain"
)

type UserUseCase interface {
	Register(ctx context.Context, email, rawPassword string) (int64, error)
	Login(ctx context.Context, email, rawPassword string) (string, string, error)
	Refresh(ctx context.Context, refreshToken string) (string, error)
	Logout(ctx context.Context, refreshToken string) error
}

type userUseCase struct {
	userRepo         domain.UserRepository
	refreshTokenRepo domain.RefreshTokenRepository
	tokenManager     domain.TokenManager
	hasher           domain.PasswordHasher
}

func NewUserUseCase(
	userRepo domain.UserRepository,
	refreshTokenRepo domain.RefreshTokenRepository,
	tokenManager domain.TokenManager,
	hasher domain.PasswordHasher,
) UserUseCase {
	return &userUseCase{
		userRepo:         userRepo,
		tokenManager:     tokenManager,
		refreshTokenRepo: refreshTokenRepo,
		hasher:           hasher,
	}
}

func (uc *userUseCase) Register(ctx context.Context, email, rawPassword string) (int64, error) {
	hashedPassword, err := uc.hasher.Hash(rawPassword)
	if err != nil {
		return 0, err
	}

	user := domain.UserWithPassword{
		User: domain.User{
			Email: email,
			Role:  roles.User,
		},
		Password: domain.HashedPassword(hashedPassword),
	}

	return uc.userRepo.CreateUser(ctx, user)
}

func (uc *userUseCase) Login(ctx context.Context, email, rawPassword string) (string, string, error) {
	user, err := uc.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return "", "", domain.ErrUserNotFound
	}

	if !uc.hasher.Compare(string(user.Password), rawPassword) {
		return "", "", domain.ErrInvalidCredentials
	}

	accessToken, err := uc.tokenManager.GenerateAccessToken(user.ID, string(user.Role))
	if err != nil {
		return "", "", err
	}

	refreshToken, err := uc.tokenManager.GenerateRefreshToken()
	if err != nil {
		return "", "", err
	}

	err = uc.refreshTokenRepo.Store(ctx, user.ID, refreshToken)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (uc *userUseCase) Refresh(ctx context.Context, refreshToken string) (string, error) {
	userID, err := uc.refreshTokenRepo.GetUserID(ctx, refreshToken)
	if err != nil {
		return "", errors.New("invalid or expired refresh token")
	}

	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return "", fmt.Errorf("user not found: %w", err)
	}

	accessToken, err := uc.tokenManager.GenerateAccessToken(user.ID, string(user.Role))
	if err != nil {
		return "", err
	}

	return accessToken, nil
}

func (uc *userUseCase) Logout(ctx context.Context, refreshToken string) error {
	err := uc.refreshTokenRepo.Delete(ctx, refreshToken)
	if err != nil {
		return errors.New("failed to delete refresh token")
	}

	return nil
}
