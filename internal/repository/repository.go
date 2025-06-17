package repository

import (
	"OrderKeeper/internal/models"
	"context"
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
