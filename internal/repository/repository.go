package repository

import (
	"OrderKeeper/internal/models"
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type Authorization interface {
	CreateUser(ctx context.Context, user models.User) (int, error)
	GetUser(ctx context.Context, username, password string) (models.User, error)
}

type Order interface {
	CreateOrder(ctx context.Context, userID int, order *models.Order) error
	GetOrders(ctx context.Context, userID int) ([]models.Order, error)
	GetOrderByID(ctx context.Context, userID int, orderID int) (models.Order, error)
	DeleteOrder(ctx context.Context, userID int, orderID int) error
	UpdateOrder(ctx context.Context, userID int, orderID int, input models.OrderUpdateInput) error
}

type Repository struct {
	Authorization
	Order
}

func NewRepository(db *pgxpool.Pool, logger *zap.Logger) *Repository {
	return &Repository{
		Authorization: NewAuthorizationRepository(db, logger),
		Order:         NewOrderRepository(db, logger),
	}
}
