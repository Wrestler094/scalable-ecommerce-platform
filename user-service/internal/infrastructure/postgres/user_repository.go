package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"

	"pkg/authenticator"

	"user-service/internal/domain"
	"user-service/internal/infrastructure/idgenerator"
	"user-service/internal/infrastructure/postgres/dao"
)

var _ domain.UserRepository = (*userRepository)(nil)

const pgErrCodeUniqueViolation = "23505"

type userRepository struct {
	router      *ShardRouter
	idGenerator idgenerator.Generator
}

func NewUserRepository(router *ShardRouter, idGenerator idgenerator.Generator) domain.UserRepository {
	return &userRepository{router: router, idGenerator: idGenerator}
}

func (r *userRepository) CreateUser(ctx context.Context, user domain.UserWithPassword) (int64, error) {
	const op = "userRepository.CreateUser"

	// Проверка существования пользователя
	existing, err := r.GetUserByEmail(ctx, user.Email)
	if err == nil && existing != nil {
		return 0, fmt.Errorf("%s: email already exists: %w", op, domain.ErrUserAlreadyExists)
	}
	if err != nil && !errors.Is(err, domain.ErrUserNotFound) {
		return 0, fmt.Errorf("%s: failed to check existing user: %w", op, err)
	}

	// Генерация ID и вставка
	id := r.idGenerator.Generate()
	db := r.router.GetShardByUserID(id)

	const query = `
        INSERT INTO users (id, email, password, role)
        VALUES ($1, $2, $3, $4)
        RETURNING id;
    `

	var idFromDB int64
	err = db.QueryRowContext(ctx, query, id, user.Email, user.Password, user.Role).Scan(&idFromDB)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == pgErrCodeUniqueViolation {
			return 0, fmt.Errorf("%s: unique violation: %w", op, domain.ErrUserAlreadyExists)
		}
		return 0, fmt.Errorf("%s: failed to insert user: %w", op, err)
	}

	return idFromDB, nil
}

func (r *userRepository) GetUserByEmail(ctx context.Context, email string) (*domain.UserWithPassword, error) {
	const op = "userRepository.GetUserByEmail"

	const query = `
		SELECT id, email, password, role, created_at, updated_at
		FROM users
		WHERE email = $1;
	`

	type result struct {
		user *domain.UserWithPassword
		err  error
	}

	resultCh := make(chan result, len(r.router.shards))
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	for i, db := range r.router.shards {
		go func(i int, db *sqlx.DB) {
			var u dao.UserModel
			err := db.GetContext(ctx, &u, query, email)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					resultCh <- result{nil, nil}
					return
				}
				resultCh <- result{nil, fmt.Errorf("%s: shard %d: %w", op, i, err)}
				return
			}

			user := &domain.UserWithPassword{
				User: domain.User{
					ID:    u.ID,
					Email: u.Email,
					Role:  authenticator.Role(u.Role),
				},
				Password: domain.HashedPassword(u.PasswordHash),
			}

			select {
			case resultCh <- result{user, nil}:
				cancel()
			case <-ctx.Done():
			}
		}(i, db)
	}

	var firstErr error
	for range r.router.shards {
		res := <-resultCh
		if res.user != nil {
			return res.user, nil
		}
		if res.err != nil && firstErr == nil {
			firstErr = res.err
		}
	}

	if firstErr != nil {
		return nil, firstErr
	}

	return nil, fmt.Errorf("%s: user not found: %w", op, domain.ErrUserNotFound)
}

func (r *userRepository) GetUserByID(ctx context.Context, id int64) (*domain.User, error) {
	const op = "userRepository.GetUserByID"

	const query = `
        SELECT id, email, password, role, created_at, updated_at
        FROM users
        WHERE id = $1;
    `

	db := r.router.GetShardByUserID(id)

	var user dao.UserModel
	if err := db.GetContext(ctx, &user, query, id); err != nil {
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
