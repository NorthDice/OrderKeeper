package service

import (
	"OrderKeeper/internal/models"
	"OrderKeeper/internal/repository"
	"context"
	"fmt"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type AuthorizationService struct {
	repo   repository.Authorization
	logger *zap.Logger
}

func NewAuthorizationService(repository repository.Authorization, logger *zap.Logger) *AuthorizationService {
	return &AuthorizationService{
		repo:   repository,
		logger: logger,
	}
}

func (a *AuthorizationService) CreateUser(ctx context.Context, user models.User) (int, error) {
	start := time.Now()

	a.logger.Info("user creation process started",
		zap.String("email", user.Email),
		zap.String("username", user.Username),
	)

	user.Password = generatePasswordHash(user.Password)
	id, err := a.repo.CreateUser(ctx, user)
	if err != nil {
		a.logger.Error("failed to create user", zap.Error(err),
			zap.String("email", user.Email),
			zap.String("username", user.Username),
			zap.Error(err),
			zap.Duration("total_duration", time.Since(start)),
		)
		return 0, fmt.Errorf("failed to create user: %w", err)
	}

	a.logger.Info("user created successfully",
		zap.Int("user_id", id),
		zap.String("email", user.Email),
		zap.String("username", user.Username),
		zap.Duration("total_service_duration", time.Since(start)),
	)

	return id, nil
}

func (a *AuthorizationService) GenerateToken(ctx context.Context, user models.User) (string, error) {
	return "", nil
}
func (a *AuthorizationService) ParseToken(ctx context.Context, token string) (int, error) {
	return 0, nil
}

func generatePasswordHash(password string) string {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		fmt.Errorf("could not generate password: %v", err)
	}

	return string(hashedPassword)
}
