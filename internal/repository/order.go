package repository

import (
	"OrderKeeper/internal/models"
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
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
