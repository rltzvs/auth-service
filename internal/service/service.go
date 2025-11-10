package service

import (
	"context"

	"auth/internal/entity"
)

type AuthRepository interface {
	CreateUser(ctx context.Context, user entity.User) error
	GetUserByEmail(ctx context.Context, email string) (entity.User, error)
}
