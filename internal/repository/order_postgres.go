package repository

import (
	"OrderKeeper/internal/models"
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"time"
)

type OrderRepository struct {
	db     *pgxpool.Pool
	logger *zap.Logger
}

func NewOrderRepository(db *pgxpool.Pool, logger *zap.Logger) *OrderRepository {
	return &OrderRepository{
		db:     db,
		logger: logger,
	}
}

func (o *OrderRepository) CreateOrder(ctx context.Context, userID int, order *models.Order) error {
	start := time.Now()
	o.logger.Debug("database insert operation started",
		zap.Int("user_id", userID),
		zap.String("status", string(order.Status)),
		zap.String("operation", "insert_order"),
	)

	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	_, err := o.db.Exec(ctx, queryInsertOrder, userID, string(order.Status))
	duration := time.Since(start)
	if err != nil {
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			o.logger.Error("database query timeout",
				zap.Int("user_id", userID),
				zap.String("status", string(order.Status)),
				zap.String("operation", "insert_order"),
				zap.Duration("timeout", time.Second*5),
				zap.Duration("actual_duration", duration),
				zap.Error(err))
			return fmt.Errorf("database query timeout: %w", err)
		}
		o.logger.Error("database insert failed",
			zap.Int("user_id", userID),
			zap.String("status", string(order.Status)),
			zap.String("operation", "insert_order"),
			zap.String("query", "INSERT INTO orders"),
			zap.Error(err),
			zap.Duration("total_duration", duration),
		)
		return fmt.Errorf("failed to create order: %w", err)
	}
	o.logger.Info("order created successfully",
		zap.Int("user_id", userID),
		zap.String("status", string(order.Status)),
		zap.Duration("total_duration", duration),
	)

	if duration > SlowQueryThreshold {
		o.logger.Warn("slow database query detected",
			zap.String("operation", "insert_order"),
			zap.Duration("db_duration", duration),
			zap.Duration("threshold", SlowQueryThreshold),
			zap.String("user_id", fmt.Sprintf("%d", userID)),
		)
	}

	return nil

}
func (o *OrderRepository) GetOrders(ctx context.Context, userID int) ([]models.Order, error) {
	return []models.Order{}, nil
}
func (o *OrderRepository) GetOrderByID(ctx context.Context, userID int, orderID int) (models.Order, error) {
	return models.Order{}, nil
}
func (o *OrderRepository) DeleteOrder(ctx context.Context, userID int, orderID int) error {
	return nil
}
func (o *OrderRepository) UpdateOrder(ctx context.Context, userID int, orderID int, input models.OrderUpdateInput) error {
	return nil
}
