package cache

import (
	"context"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

type RedisCache struct {
	Client *redis.Client
	logger *zap.Logger
}

type RedisConfig struct {
	Address  string
	Password string
	Database int
}

func NewRedisCache(ctx context.Context, cfg RedisConfig, logger *zap.Logger) (RedisCache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Address,
		Password: cfg.Password,
		DB:       cfg.Database,
	})

	return RedisCache{
		Client: client,
		logger: logger,
	}, nil
}
