package postgres

import (
	"OrderKeeper/internal/handler/metrics"
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

type CachedOrderRepository struct {
	orderRepo *OrderRepository
	cache     *cache.RedisCache
	logger    *zap.Logger
}

func NewCachedOrderRepository(db *pgxpool.Pool, cache *cache.RedisCache, logger *zap.Logger) *CachedOrderRepository {
	return &CachedOrderRepository{
		orderRepo: NewOrderRepository(db, logger),
		cache:     cache,
		logger:    logger,
	}
}
func (c *CachedOrderRepository) CreateOrder(ctx context.Context, userID int, order *models.Order) error {
	err := c.orderRepo.CreateOrder(ctx, userID, order)
	if err != nil {
		c.logger.Error("failed to create order in order repository", zap.Error(err))
		return err
	}

	metrics.RecordOrder("created")

	listCacheKey := fmt.Sprintf("orders:user:%d", userID)
	if cacheErr := c.cache.Delete(ctx, listCacheKey); cacheErr != nil {
		c.logger.Warn("Failed to invalidate order list cache",
			zap.Error(cacheErr),
			zap.Int("user_id", userID),
		)
	}

	orderCacheKey := fmt.Sprintf("order:user:%d:id:%d", userID, order.ID)
	if cacheErr := c.cache.Set(ctx, orderCacheKey, order, 30*time.Minute); cacheErr != nil {
		c.logger.Warn("Failed to cache created order",
			zap.Error(cacheErr),
			zap.Int("user_id", userID),
			zap.Int("order_id", order.ID),
		)
	}

	return nil
}
func (c *CachedOrderRepository) GetOrders(ctx context.Context, userID int) ([]models.Order, error) {
	cacheKey := fmt.Sprintf("orders:user:%d", userID)

	var cachedOrders []models.Order
	err := c.cache.Get(ctx, cacheKey, &cachedOrders)
	if err == nil {
		c.logger.Debug("Orders found in cache", zap.Int("userID", userID))
		metrics.RecordCacheHit("user_orders")
		return cachedOrders, nil
	}

	if !errors.Is(err, redis.Nil) {
		c.logger.Warn("Redis error when getting orders",
			zap.Error(err),
			zap.Int("userID", userID),
		)
	}

	metrics.RecordCacheMiss("user_orders")

	orders, err := c.orderRepo.GetOrders(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get orders from orders repository: %w", err)
	}

	if cacheErr := c.cache.Set(ctx, cacheKey, orders, 15*time.Minute); cacheErr != nil {
		c.logger.Warn("Failed to cache orders",
			zap.Error(cacheErr),
			zap.Int("userID", userID))
	}

	c.logger.Debug("Orders loaded from database and cached", zap.Int("userID", userID))
	return orders, nil
}

func (c *CachedOrderRepository) GetOrderByID(ctx context.Context, userID int, orderID int) (models.Order, error) {
	cacheKey := fmt.Sprintf("order:user:%d:id:%d", userID, orderID)
	var cachedOrder models.Order

	err := c.cache.Get(ctx, cacheKey, &cachedOrder)
	if err == nil {
		c.logger.Debug("Order found in cache",
			zap.Int("userID", userID),
			zap.Int("orderID", orderID),
		)
		metrics.RecordCacheHit("order")
		return cachedOrder, nil
	}

	if !errors.Is(err, redis.Nil) {
		c.logger.Warn("Redis error when getting order",
			zap.Error(err),
			zap.Int("userID", userID),
			zap.Int("orderID", orderID))
	}

	metrics.RecordCacheMiss("order")

	order, err := c.orderRepo.GetOrderByID(ctx, userID, orderID)
	if err != nil {
		return models.Order{}, err
	}

	if cacheErr := c.cache.Set(ctx, cacheKey, order, 30*time.Minute); cacheErr != nil {
		c.logger.Warn("Failed to cache order",
			zap.Error(cacheErr),
			zap.Int("orderID", orderID))
	}

	c.logger.Debug("Order loaded from database and cached",
		zap.Int("userID", userID),
		zap.Int("orderID", orderID))
	return order, nil
}
func (c *CachedOrderRepository) UpdateOrder(ctx context.Context, userID int, orderID int, input models.OrderUpdateInput) error {

	err := c.orderRepo.UpdateOrder(ctx, userID, orderID, input)
	if err != nil {
		return err
	}

	metrics.RecordOrder("updated")

	orderCacheKey := fmt.Sprintf("order:user:%d:id:%d", userID, orderID)
	listCacheKey := fmt.Sprintf("orders:user:%d", userID)

	if cacheErr := c.cache.Delete(ctx, orderCacheKey); cacheErr != nil {
		c.logger.Warn("Failed to invalidate order cache",
			zap.Error(cacheErr),
			zap.Int("orderID", orderID))
	}

	if cacheErr := c.cache.Delete(ctx, listCacheKey); cacheErr != nil {
		c.logger.Warn("Failed to invalidate orders list cache",
			zap.Error(cacheErr),
			zap.Int("userID", userID))
	}

	return nil
}

func (c *CachedOrderRepository) DeleteOrder(ctx context.Context, userID int, orderID int) error {

	err := c.orderRepo.DeleteOrder(ctx, userID, orderID)
	if err != nil {
		return err
	}

	metrics.RecordOrder("deleted")

	orderCacheKey := fmt.Sprintf("order:user:%d:id:%d", userID, orderID)
	listCacheKey := fmt.Sprintf("orders:user:%d", userID)

	if cacheErr := c.cache.Delete(ctx, orderCacheKey); cacheErr != nil {
		c.logger.Warn("Failed to delete order from cache",
			zap.Error(cacheErr),
			zap.Int("orderID", orderID))
	}

	if cacheErr := c.cache.Delete(ctx, listCacheKey); cacheErr != nil {
		c.logger.Warn("Failed to invalidate orders list cache",
			zap.Error(cacheErr),
			zap.Int("userID", userID))
	}

	return nil
}
