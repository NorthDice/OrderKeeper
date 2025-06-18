package service

import (
	"OrderKeeper/internal/models"
	"OrderKeeper/internal/repository"
	"context"
	"go.uber.org/zap"
)

type AuthorizationService struct {
	repository.Authorization
	logger *zap.Logger
}

func NewAuthorizationService(repository repository.Authorization, logger *zap.Logger) *AuthorizationService {
	return &AuthorizationService{
		Authorization: repository,
		logger:        logger,
	}
}

func (a *AuthorizationService) CreateUser(ctx context.Context, user models.User) (int, error) {
	return 0, nil
}
func (a *AuthorizationService) GenerateToken(ctx context.Context, user models.User) (string, error) {
	return "", nil
}
func (a *AuthorizationService) ParseToken(ctx context.Context, token string) (int, error) {
	return 0, nil
}
