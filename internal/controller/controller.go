package controller

import (
	"context"

	"auth/internal/entity"
)

type AuthService interface {
	Login(ctx context.Context, user entity.User) (string, error)
	Register(ctx context.Context, user entity.User) (string, error)
}
