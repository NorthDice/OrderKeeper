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
	return nil, nil
}
func (o *OrderService) GetOrderByID(ctx context.Context, userID int, orderID int) (models.Order, error) {
	return models.Order{}, nil
}
func (o *OrderService) DeleteOrder(ctx context.Context, userID int, orderID int) error {
	return nil
}
func (o *OrderService) UpdateOrder(ctx context.Context, userID int, orderID int, input models.OrderUpdateInput) error {
	return nil
}
