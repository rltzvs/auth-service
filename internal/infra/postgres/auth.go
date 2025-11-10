package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"auth/internal/entity"
)

type AuthRepository struct {
	Pool *pgxpool.Pool
}

func NewAuthRepository(pool *pgxpool.Pool) *AuthRepository {
	return &AuthRepository{
		Pool: pool,
	}
}

func (r *AuthRepository) GetUserByEmail(ctx context.Context, email string) (entity.User, error) {
	query := `
			SELECT id, email, password
			FROM users
			WHERE email = $1
	`

	var user entity.User
	err := r.Pool.QueryRow(ctx, query, email).Scan(&user.ID, &user.Email, &user.Password)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.User{}, entity.ErrUserNotFound
		}
		return entity.User{}, fmt.Errorf("failed to get user by email: %w", err)
	}

	return user, nil
}

func (r *AuthRepository) CreateUser(ctx context.Context, user entity.User) error {
	query := `
			INSERT INTO users (email, password)
			VALUES ($1, $2)
	`

	_, err := r.Pool.Exec(ctx, query, user.Email, user.Password)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return entity.ErrUserAlreadyExists
		}
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}
