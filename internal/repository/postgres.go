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
const (
	queryInsertUser = `
		INSERT INTO users (username, email, password)
		VALUES ($1, $2, $3)
		RETURNING id
	`
	querySelectUser = `
		SELECT id, username, email, password FROM users
		WHERE username = $1
	`
)
const (
	queryInsertOrder = `
		INSERT INTO orders (user_id, status)
	    VALUES ($1, $2)
	`
	querySelectOrdersByUser = `
	SELECT id, user_id, status,  created_at, updated_at
	FROM orders
	WHERE user_id = $1
	`
	querySelectOrderByID = `
		SELECT id, user_id, status, created_at, updated_at
		FROM orders
		WHERE user_id = $1 AND id = $2
	   `
	queryDeleteOrderByID = `
		DELETE FROM orders
		WHERE user_id = $1 AND id = $2
	`
	queryUpdateOrderByID = `
		UPDATE orders
		SET status = $1, updated_at = NOW()
		WHERE user_id = $2 AND id = $3
		`
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
