package service

import (
	"OrderKeeper/internal/models"
	"OrderKeeper/internal/repository"
	"context"
	"fmt"
	"go.uber.org/zap"
	"time"
)

type OrderService struct {
	repository repository.Order
	logger     *zap.Logger
}

func NewOrderService(repo repository.Order, logger *zap.Logger) *OrderService {
	return &OrderService{
		repository: repo,
		logger:     logger,
	}
}

func (o *OrderService) CreateOrder(ctx context.Context, userID int, order *models.Order) error {

	start := time.Now()
	o.logger.Info("order creation process started",
		zap.Int("user_id", userID),
		zap.String("status", string(order.Status)),
	)

	err := o.repository.CreateOrder(ctx, userID, order)
	if err != nil {
		o.logger.Error("failed to create order",
			zap.Int("user_id", userID),
			zap.String("status", string(order.Status)),
			zap.Error(err),
			zap.Duration("total_duration", time.Since(start)),
		)
		return fmt.Errorf("failed to create order: %w", err)
	}

	o.logger.Info("order created successfully",
		zap.Int("user_id", userID),
		zap.String("status", string(order.Status)),
		zap.Duration("total_service_duration", time.Since(start)),
	)

	return nil
}
func (o *OrderService) GetOrders(ctx context.Context, userID int) ([]models.Order, error) {
	start := time.Now()
	o.logger.Info("fetching orders for user",
		zap.Int("user_id", userID),
	)
	orders, err := o.repository.GetOrders(ctx, userID)
	if err != nil {
		o.logger.Error("failed to fetch orders",
			zap.Int("user_id", userID),
			zap.Error(err),
			zap.Duration("total_duration", time.Since(start)),
		)
		return nil, fmt.Errorf("failed to fetch orders: %w", err)
	}
	o.logger.Info("orders fetched successfully",
		zap.Int("user_id", userID),
		zap.Int("order_count", len(orders)),
		zap.Duration("total_duration", time.Since(start)),
	)
	return orders, nil
}
func (o *OrderService) GetOrderByID(ctx context.Context, userID int, orderID int) (models.Order, error) {
	start := time.Now()
	o.logger.Info("fetching order by ID",
		zap.Int("user_id", userID),
		zap.Int("order_id", orderID),
	)

	order, err := o.repository.GetOrderByID(ctx, userID, orderID)
	if err != nil {
		o.logger.Error("failed to fetch order by ID",
			zap.Int("user_id", userID),
			zap.Int("order_id", orderID),
			zap.Error(err),
			zap.Duration("total_duration", time.Since(start)),
		)
		return models.Order{}, fmt.Errorf("failed to fetch order by ID: %w", err)
	}
	o.logger.Info("order fetched successfully",
		zap.Int("user_id", userID),
		zap.Int("order_id", orderID),
		zap.Duration("total_duration", time.Since(start)),
	)
	return order, nil
}
func (o *OrderService) DeleteOrder(ctx context.Context, userID int, orderID int) error {
	start := time.Now()
	o.logger.Info("deleting order",
		zap.Int("user_id", userID),
		zap.Int("order_id", orderID),
	)

	err := o.repository.DeleteOrder(ctx, userID, orderID)
	if err != nil {
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
	return nil
}
func (o *OrderService) UpdateOrder(ctx context.Context, userID int, orderID int, input models.OrderUpdateInput) error {
	start := time.Now()
	o.logger.Info("updating order",
		zap.Int("user_id", userID),
		zap.Int("order_id", orderID),
	)
	err := o.repository.UpdateOrder(ctx, userID, orderID, input)
	if err != nil {
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
	return nil
}
