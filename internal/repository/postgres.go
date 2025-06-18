package repository

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	userTable   = "users"
	ordersTable = "orders"
)

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

func NewPostgresDB(ctx context.Context, cfg Config) (*pgxpool.Pool, error) {
	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, cfg.SSLMode)
	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, err
	}
	err = pool.Ping(ctx)
	if err != nil {
		pool.Close()
		return nil, err
	}
	return pool, nil
}
