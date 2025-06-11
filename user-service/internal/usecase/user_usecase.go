package usecase

import (
	"context"
	"errors"
	"fmt"

	"pkg/authenticator"
	"user-service/internal/domain"
)

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
) domain.UserUseCase {
	return &userUseCase{
		userRepo:         userRepo,
		tokenManager:     tokenManager,
		refreshTokenRepo: refreshTokenRepo,
		hasher:           hasher,
	}
}

func (uc *userUseCase) Register(ctx context.Context, email, rawPassword string) (int64, error) {
	const op = "userUseCase.Register"

	hashedPassword, err := uc.hasher.Hash(rawPassword)
	if err != nil {
		return 0, fmt.Errorf("%s: failed to hash password: %w", op, err)
	}

	user := domain.UserWithPassword{
		User: domain.User{
			Email: email,
			Role:  authenticator.User,
		},
		Password: domain.HashedPassword(hashedPassword),
	}

	id, err := uc.userRepo.CreateUser(ctx, user)
	if err != nil {
		if errors.Is(err, domain.ErrUserAlreadyExists) {
			return 0, domain.ErrUserAlreadyExists
		}
		return 0, fmt.Errorf("%s: failed to create user: %w", op, err)
	}

	return id, nil
}

func (uc *userUseCase) Login(ctx context.Context, email, rawPassword string) (string, string, error) {
	const op = "userUseCase.Login"

	user, err := uc.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return "", "", domain.ErrInvalidCredentials
		}
		return "", "", fmt.Errorf("%s: failed to get user: %w", op, err)
	}

	isCorrect, err := uc.hasher.Compare(string(user.Password), rawPassword)
	if err != nil {
		return "", "", fmt.Errorf("%s: failed to compare password: %w", op, err)
	}
	if !isCorrect {
		return "", "", domain.ErrInvalidCredentials
	}

	accessToken, err := uc.tokenManager.GenerateAccessToken(user.ID, string(user.Role))
	if err != nil {
		return "", "", fmt.Errorf("%s: failed to generate access token: %w", op, err)
	}

	refreshToken, err := uc.tokenManager.GenerateRefreshToken()
	if err != nil {
		return "", "", fmt.Errorf("%s: failed to generate refresh token: %w", op, err)
	}

	err = uc.refreshTokenRepo.Store(ctx, user.ID, refreshToken)
	if err != nil {
		return "", "", fmt.Errorf("%s: failed to store refresh token: %w", op, err)
	}

	return accessToken, refreshToken, nil
}

func (uc *userUseCase) Refresh(ctx context.Context, refreshToken string) (string, error) {
	const op = "userUseCase.Refresh"

	userID, err := uc.refreshTokenRepo.GetUserID(ctx, refreshToken)
	if err != nil {
		return "", fmt.Errorf("%s: invalid or expired refresh token: %w", op, err)
	}

	user, err := uc.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return "", domain.ErrUserNotFound
		}
		return "", fmt.Errorf("%s: failed to get user: %w", op, err)
	}

	accessToken, err := uc.tokenManager.GenerateAccessToken(user.ID, string(user.Role))
	if err != nil {
		return "", fmt.Errorf("%s: failed to generate access token: %w", op, err)
	}

	return accessToken, nil
}

func (uc *userUseCase) Logout(ctx context.Context, refreshToken string) error {
	const op = "userUseCase.Logout"

	err := uc.refreshTokenRepo.Delete(ctx, refreshToken)
	if err != nil {
		return fmt.Errorf("%s: failed to delete refresh token: %w", op, err)
	}

	return nil
}
