package repository

import (
	"OrderKeeper/internal/models"
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type AuthorizationRepository struct {
	db     *pgxpool.Pool
	logger *zap.Logger
}

func NewAuthorizationRepository(db *pgxpool.Pool, logger *zap.Logger) *AuthorizationRepository {
	return &AuthorizationRepository{
		db:     db,
		logger: logger,
	}
}

func (a *AuthorizationRepository) CreateUser(ctx context.Context, user models.User) (int, error) {
	return 0, nil
}

func (a *AuthorizationRepository) GetUser(ctx context.Context, username, password string) (models.User, error) {
	return models.User{}, nil
}
