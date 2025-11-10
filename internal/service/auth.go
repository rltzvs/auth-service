package service

import (
	"context"

	"auth/internal/entity"
	"auth/internal/security"
)

type AuthService struct {
	repo           AuthRepository
	passwordHasher security.PasswordHasher
	jwtManager     security.JWTManager
}

func NewAuthService(repo AuthRepository, passwordHasher security.PasswordHasher, jwtManager security.JWTManager) *AuthService {
	return &AuthService{repo: repo, passwordHasher: passwordHasher, jwtManager: jwtManager}
}

func (s *AuthService) Register(ctx context.Context, user entity.User) (string, error) {
	pass, err := s.passwordHasher.HashPassword(user.Password)
	if err != nil {
		return "", err
	}
	user.Password = pass

	err = s.repo.CreateUser(ctx, user)
	if err != nil {
		return "", err
	}
	return "", nil
}

func (s *AuthService) Login(ctx context.Context, userData entity.User) (string, error) {
	user, err := s.repo.GetUserByEmail(ctx, userData.Email)
	if err != nil {
		return "", err
	}

	err = s.passwordHasher.CheckPasswordHash(userData.Password, user.Password)
	if err != nil {
		return "", err
	}

	token, err := s.jwtManager.GenerateToken(int64(user.ID))
	if err != nil {
		return "", err
	}
	return token, nil
}
