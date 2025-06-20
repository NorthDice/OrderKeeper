package repository

import (
	"OrderKeeper/internal/models"
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
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
	start := time.Now()
	o.logger.Debug("fetching orders for user",
		zap.Int("user_id", userID),
		zap.String("operation", "get_orders"),
	)

	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()
	rows, err := o.db.Query(ctx, querySelectOrdersByUser, userID)
	if err != nil {
		o.logger.Error("failed to fetch orders",
			zap.Int("user_id", userID),
			zap.Error(err),
			zap.Duration("total_duration", time.Since(start)),
		)
		return nil, fmt.Errorf("failed to fetch orders: %w", err)
	}
	defer rows.Close()
	duration := time.Since(start)

	var orders []models.Order
	for rows.Next() {
		var order models.Order
		err := rows.Scan(&order.ID, &order.UserID, &order.Status, &order.CreatedAt, &order.UpdatedAt)
		if err != nil {
			o.logger.Error("failed to scan order",
				zap.Int("user_id", userID),
				zap.Error(err),
				zap.Duration("total_duration", time.Since(start)),
			)
			return nil, fmt.Errorf("failed to scan order: %w", err)
		}
		orders = append(orders, order)
	}

	o.logger.Info("orders fetched successfully",
		zap.Int("user_id", userID),
		zap.Int("order_count", len(orders)),
		zap.Duration("total_duration", time.Since(start)),
	)

	if duration > SlowQueryThreshold {
		o.logger.Warn("slow database query detected",
			zap.String("operation", "get_orders"),
			zap.Duration("db_duration", duration),
			zap.Duration("threshold", SlowQueryThreshold),
			zap.Int("user_id", userID),
		)
	}

	return orders, nil
}
func (o *OrderRepository) GetOrderByID(ctx context.Context, userID int, orderID int) (models.Order, error) {
	start := time.Now()

	o.logger.Debug("fetching order by ID",
		zap.Int("user_id", userID),
		zap.Int("order_id", orderID),
		zap.String("operation", "get_order_by_id"),
	)

	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()
	row := o.db.QueryRow(ctx, querySelectOrderByID, userID, orderID)
	duration := time.Since(start)

	var order models.Order
	err := row.Scan(&order.ID, &order.UserID, &order.Status, &order.CreatedAt, &order.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			o.logger.Error("order not found",
				zap.Int("user_id", userID),
				zap.Int("order_id", orderID),
				zap.Error(err),
				zap.Duration("total_duration", duration),
			)
			return models.Order{}, fmt.Errorf("order not found: %w", err)
		}
		o.logger.Error("failed to fetch order",
			zap.Int("user_id", userID),
			zap.Int("order_id", orderID),
			zap.Error(err),
			zap.Duration("total_duration", duration),
		)
		return models.Order{}, fmt.Errorf("failed to fetch order: %w", err)
	}

	o.logger.Info("order fetched successfully",
		zap.Int("user_id", userID),
		zap.Int("order_id", orderID),
		zap.Duration("total_duration", duration),
	)
	if duration > SlowQueryThreshold {
		o.logger.Warn("slow database query detected",
			zap.String("operation", "get_order_by_id"),
			zap.Duration("db_duration", duration),
			zap.Duration("threshold", SlowQueryThreshold),
			zap.Int("user_id", userID),
			zap.Int("order_id", orderID),
		)
	}

	return order, nil
}
func (o *OrderRepository) DeleteOrder(ctx context.Context, userID int, orderID int) error {
	start := time.Now()
	o.logger.Debug("deleting order",
		zap.Int("user_id", userID),
		zap.Int("order_id", orderID),
		zap.String("operation", "delete_order"),
	)
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()
	_, err := o.db.Exec(ctx, queryDeleteOrderByID, userID, orderID)
	duration := time.Since(start)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			o.logger.Error("order not found for deletion",
				zap.Int("user_id", userID),
				zap.Int("order_id", orderID),
				zap.Error(err),
				zap.Duration("total_duration", time.Since(start)),
			)
			return fmt.Errorf("order not found: %w", err)
		}

		o.logger.Error("failed to delete order",
			zap.Int("user_id", userID),
			zap.Int("order_id", orderID),
			zap.Error(err),
			zap.Duration("total_duration", time.Since(start)),
		)
		return fmt.Errorf("failed to delete order: %w", err)
	}
	o.logger.Info("order deleted successfully",
		zap.Int("user_id", userID),
		zap.Int("order_id", orderID),
		zap.Duration("total_duration", time.Since(start)),
	)
	if duration > SlowQueryThreshold {
		o.logger.Warn("slow database query detected",
			zap.String("operation", "delete_order"),
			zap.Duration("db_duration", duration),
			zap.Duration("threshold", SlowQueryThreshold),
			zap.Int("user_id", userID),
			zap.Int("order_id", orderID),
		)
	}
	return nil
}
func (o *OrderRepository) UpdateOrder(ctx context.Context, userID int, orderID int, input models.OrderUpdateInput) error {
	start := time.Now()
	o.logger.Debug("updating order",
		zap.Int("user_id", userID),
		zap.Int("order_id", orderID),
		zap.String("operation", "update_order"),
	)
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()
	_, err := o.db.Exec(ctx, queryUpdateOrderByID, input.Status, userID, orderID)
	duration := time.Since(start)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			o.logger.Error("order not found for update",
				zap.Int("user_id", userID),
				zap.Int("order_id", orderID),
				zap.Error(err),
				zap.Duration("total_duration", time.Since(start)),
			)
			return fmt.Errorf("order not found: %w", err)
		}
		o.logger.Error("failed to update order",
			zap.Int("user_id", userID),
			zap.Int("order_id", orderID),
			zap.Error(err),
			zap.Duration("total_duration", time.Since(start)),
		)
		return fmt.Errorf("failed to update order: %w", err)
	}
	o.logger.Info("order updated successfully",
		zap.Int("user_id", userID),
		zap.Int("order_id", orderID),
		zap.Duration("total_duration", time.Since(start)),
	)
	if duration > SlowQueryThreshold {
		o.logger.Warn("slow database query detected",
			zap.String("operation", "update_order"),
			zap.Duration("db_duration", duration),
			zap.Duration("threshold", SlowQueryThreshold),
			zap.Int("user_id", userID),
			zap.Int("order_id", orderID),
		)
	}
	return nil
}
