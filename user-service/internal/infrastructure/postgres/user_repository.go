package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"pkg/authenticator"
	"user-service/internal/domain"
	"user-service/internal/infrastructure/postgres/dao"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

var _ domain.UserRepository = (*userRepository)(nil)

type userRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) domain.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) CreateUser(ctx context.Context, user domain.UserWithPassword) (int64, error) {
	const op = "userRepository.CreateUser"

	const query = `
        INSERT INTO users (email, password, role)
        VALUES ($1, $2, $3)
        RETURNING id;
    `

	var id int64
	err := r.db.QueryRowContext(ctx, query, user.Email, user.Password, user.Role).Scan(&id)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return 0, fmt.Errorf("%s: failed to insert user: %w", op, domain.ErrUserAlreadyExists)
		}
		return 0, fmt.Errorf("%s: failed to insert user: %w", op, err)
	}

	return id, nil
}

func (r *userRepository) GetUserByEmail(ctx context.Context, email string) (*domain.UserWithPassword, error) {
	const op = "userRepository.GetUserByEmail"

	const query = `
        SELECT id, email, password, role, created_at, updated_at
        FROM users
        WHERE email = $1;
    `

	var user dao.UserModel
	if err := r.db.GetContext(ctx, &user, query, email); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%s: failed to get user by email: %w", op, domain.ErrUserNotFound)
		}
		return nil, fmt.Errorf("%s: failed to get user by email: %w", op, err)
	}

	userWithPassword := domain.UserWithPassword{
		User: domain.User{
			ID:    user.ID,
			Email: user.Email,
			Role:  authenticator.Role(user.Role),
		},
		Password: domain.HashedPassword(user.PasswordHash),
	}

	return &userWithPassword, nil
}

func (r *userRepository) GetUserByID(ctx context.Context, id int64) (*domain.User, error) {
	const op = "userRepository.GetUserByID"

	const query = `
        SELECT id, email, password, role, created_at, updated_at
        FROM users
        WHERE id = $1;
    `

	var user dao.UserModel
	if err := r.db.GetContext(ctx, &user, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%s: failed to get user by ID: %w", op, domain.ErrUserNotFound)
		}
		return nil, fmt.Errorf("%s: failed to get user by ID: %w", op, err)
	}

	outputUser := domain.User{
		ID:    user.ID,
		Email: user.Email,
		Role:  authenticator.Role(user.Role),
	}

	return &outputUser, nil
}
