package postgres

import (
	"OrderKeeper/internal/models"
	"OrderKeeper/internal/repository/cache"
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"time"
)

type CachedAuthRepository struct {
	authRepo *AuthorizationRepository
	cache    *cache.RedisCache
	logger   *zap.Logger
}

func NewCachedAuthRepository(db *pgxpool.Pool, cache *cache.RedisCache, logger *zap.Logger) *CachedAuthRepository {
	return &CachedAuthRepository{
		authRepo: NewAuthorizationRepository(db, logger),
		cache:    cache,
		logger:   logger,
	}
}
func (c *CachedAuthRepository) CreateUser(ctx context.Context, user models.User) (int, error) {
	userID, err := c.authRepo.CreateUser(ctx, user)
	if err != nil {
		c.logger.Error("failed to create user in auth repository", zap.Error(err))
		return 0, fmt.Errorf("failed to create user in auth repository: %w", err)
	}

	user.ID = userID
	cacheKey := fmt.Sprintf("user:username:%s", user.Username)
	if cacheErr := c.cache.Set(ctx, cacheKey, user, 1*time.Hour); cacheErr != nil {
		c.logger.Warn("Failed to cache created user",
			zap.Error(cacheErr),
			zap.String("username", user.Username),
		)
	}

	return userID, nil
}

func (c *CachedAuthRepository) GetUser(ctx context.Context, username, password string) (models.User, error) {
	cacheKey := fmt.Sprintf("user:username:%s", username)

	var cachedUser models.User
	err := c.cache.Get(ctx, cacheKey, &cachedUser)
	if err == nil {
		if cachedUser.Password == password {
			c.logger.Debug("User retrieved from cache",
				zap.String("username", username),
			)
			return cachedUser, nil
		}
		return models.User{}, fmt.Errorf("invalid credentials")
	}

	if errors.Is(err, redis.Nil) {
		c.logger.Debug("User not found in cache, querying database",
			zap.String("username", username),
		)
	}
	user, err := c.authRepo.GetUser(ctx, username, password)
	if err != nil {
		c.logger.Error("failed to get user from auth repository",
			zap.Error(err),
			zap.String("username", username),
		)
		return models.User{}, fmt.Errorf("failed to get user from auth repository: %w", err)
	}

	if cacheErr := c.cache.Set(ctx, cacheKey, user, 1*time.Hour); cacheErr != nil {
		c.logger.Warn("Failed to cache retrieved user",
			zap.Error(cacheErr),
			zap.String("username", username),
		)
	}

	c.logger.Debug("User retrieved from database and cached",
		zap.String("username", username),
	)

	return user, nil
}
