package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"time"
)

func (r *RedisCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		r.logger.Error("failed to marshal value", zap.Error(err))
		return fmt.Errorf("failed to marshal: %w", err)
	}
	err = r.Client.Set(ctx, key, data, ttl).Err()
	if err != nil {
		r.logger.Error("failed to set value in cache", zap.String("key", key), zap.Error(err))
		return fmt.Errorf("failed to set value in cache: %w", err)
	}
	r.logger.Info("value set in cache", zap.String("key", key), zap.Duration("ttl", ttl))
	return nil
}
func (r *RedisCache) Get(ctx context.Context, key string, dest interface{}) error {
	val, err := r.Client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			r.logger.Info("key not found in cache", zap.String("key", key))
			return nil // Key not found, return nil
		}
		r.logger.Error("failed to get value from cache", zap.String("key", key), zap.Error(err))
		return fmt.Errorf("failed to get value from cache: %w", err)
	}
	err = json.Unmarshal([]byte(val), dest)
	if err != nil {
		r.logger.Error("failed to unmarshal value from cache", zap.String("key", key), zap.Error(err))
		return fmt.Errorf("failed to unmarshal value: %w", err)
	}
	r.logger.Info("value retrieved from cache", zap.String("key", key))
	return nil
}
func (r *RedisCache) Delete(ctx context.Context, key string) error {
	err := r.Client.Del(ctx, key).Err()
	if err != nil {
		r.logger.Error("failed to delete key from cache", zap.String("key", key), zap.Error(err))
		return fmt.Errorf("failed to delete key from cache: %w", err)
	}
	r.logger.Info("key deleted from cache", zap.String("key", key))
	return nil
}
func (r *RedisCache) DeletePattern(ctx context.Context, pattern string) error {
	keys, err := r.Client.Keys(ctx, pattern).Result()
	if err != nil {
		r.logger.Error("failed to get keys by pattern", zap.String("pattern", pattern), zap.Error(err))
		return fmt.Errorf("failed to get keys by pattern: %w", err)
	}
	if len(keys) > 0 {
		return r.Client.Del(ctx, keys...).Err()
	}
	r.logger.Info("no keys found for pattern", zap.String("pattern", pattern))
	return nil
}

func (r *RedisCache) Ping(ctx context.Context) error {
	return r.Client.Ping(ctx).Err()
}

func (r *RedisCache) Close() error {
	return r.Client.Close()
}
