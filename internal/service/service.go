package service

import (
	"OrderKeeper/internal/models"
	"OrderKeeper/internal/repository"
	"context"
	"go.uber.org/zap"
)

type Authorization interface {
	CreateUser(ctx context.Context, user models.User) (int, error)
	GenerateToken(ctx context.Context, user models.User) (string, error)
	ParseToken(ctx context.Context, token string) (int, error)
}

type Order interface {
	CreateOrder(ctx context.Context, userID int, order *models.Order) error
	GetOrders(ctx context.Context, userID int) ([]models.Order, error)
	GetOrderByID(ctx context.Context, userID int, orderID int) (models.Order, error)
	DeleteOrder(ctx context.Context, userID int, orderID int) error
	UpdateOrder(ctx context.Context, userID int, orderID int, input models.OrderUpdateInput) error
}

type Service struct {
	Authorization
	Order
}

func NewService(repo repository.Repository, logger zap.Logger) *Service {
	return &Service{
		Authorization: NewAuthorizationService(repo.Authorization, logger),
		Order:         NewOrderService(repo.Order, logger),
	}
}
