package service

import (
	"OrderKeeper/internal/models"
	"OrderKeeper/internal/repository"
	"context"
	"go.uber.org/zap"
)

type OrderService struct {
	repository repository.Order
	logger     zap.Logger
}

func NewOrderService(repo repository.Order, logger zap.Logger) *OrderService {
	return &OrderService{
		repository: repo,
		logger:     logger,
	}
}

func (o *OrderService) CreateOrder(ctx context.Context, userID int, order *models.Order) error {
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
