package postgres

import (
	"context"

	"pkg/roles"
	"user-service/internal/domain"
	"user-service/internal/infrastructure/postgres/dao"

	"github.com/jmoiron/sqlx"
)

var _ domain.UserRepository = (*userRepository)(nil)

type userRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) domain.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) CreateUser(ctx context.Context, user domain.UserWithPassword) (int64, error) {
	query := `
        INSERT INTO users (email, password, role)
        VALUES ($1, $2, $3)
        RETURNING id;
    `

	var id int64
	err := r.db.QueryRowContext(ctx, query, user.Email, user.Password, user.Role).Scan(&id)

	return id, err
}

func (r *userRepository) GetUserByEmail(ctx context.Context, email string) (*domain.UserWithPassword, error) {
	const query = `
        SELECT id, email, password, role, created_at, updated_at
        FROM users
        WHERE email = $1;
    `

	var user dao.UserModel
	err := r.db.GetContext(ctx, &user, query, email)
	if err != nil {
		return nil, err
	}

	userWithPassword := domain.UserWithPassword{
		User: domain.User{
			ID:    user.ID,
			Email: user.Email,
			Role:  roles.Role(user.Role),
		},
		Password: domain.HashedPassword(user.PasswordHash),
	}

	return &userWithPassword, nil
}

func (r *userRepository) GetByID(ctx context.Context, id int64) (*domain.User, error) {
	const query = `
        SELECT id, email, password, role, created_at, updated_at
        FROM users
        WHERE id = $1;
    `

	var user dao.UserModel
	if err := r.db.GetContext(ctx, &user, query, id); err != nil {
		return nil, err
	}

	outputUser := domain.User{
		ID:    user.ID,
		Email: user.Email,
		Role:  roles.Role(user.Role),
	}

	return &outputUser, nil
}
