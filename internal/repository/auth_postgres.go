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

const (
	DefaultDBTimeout   = 5 * time.Second
	SlowQueryThreshold = 100 * time.Millisecond
)

type AuthorizationRepository struct {
	db     *pgxpool.Pool
	logger *zap.Logger
}

func NewAuthorizationRepository(db *pgxpool.Pool, logger *zap.Logger) *AuthorizationRepository {
	return &AuthorizationRepository{
		db:     db,
		logger: logger,
	}
}

func (a *AuthorizationRepository) CreateUser(ctx context.Context, user models.User) (int, error) {
	start := time.Now()

	a.logger.Debug("database insert operation started",
		zap.String("email", user.Email),
		zap.String("username", user.Username),
		zap.String("operation", "insert_user"),
	)

	ctx, cancel := context.WithTimeout(ctx, DefaultDBTimeout)
	defer cancel()

	query := `INSERT INTO users (username, email, password)
              VALUES ($1, $2, $3)
              RETURNING id`

	var id int
	err := a.db.QueryRow(ctx, query, user.Username, user.Email, user.Password).Scan(&id)
	duration := time.Since(start)
	if err != nil {
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			a.logger.Error("database query timeout",
				zap.String("email", user.Email),
				zap.String("username", user.Username),
				zap.String("operation", "insert_user"),
				zap.Duration("timeout", DefaultDBTimeout),
				zap.Duration("actual_duration", duration),
				zap.Error(err))
			return 0, fmt.Errorf("database query timeout after %v: %w", DefaultDBTimeout, err)
		}

		a.logger.Error("database insert failed",
			zap.String("email", user.Email),
			zap.String("username", user.Username),
			zap.String("operation", "insert_user"),
			zap.String("query", "INSERT INTO users"),
			zap.Error(err),
			zap.Duration("db_duration", duration),
		)
		return 0, fmt.Errorf("could not create user: %w", err)
	}

	a.logger.Info("user inserted successfully",
		zap.Int("user_id", id),
		zap.String("email", user.Email),
		zap.String("username", user.Username),
		zap.String("operation", "insert_user"),
		zap.Duration("db_duration", duration),
	)

	if duration > SlowQueryThreshold {
		a.logger.Warn("slow database query detected",
			zap.String("operation", "insert_user"),
			zap.Duration("db_duration", duration),
			zap.Duration("threshold", SlowQueryThreshold),
			zap.String("email", user.Email))
	}

	return id, nil
}

func (a *AuthorizationRepository) GetUser(ctx context.Context, username, password string) (models.User, error) {
	start := time.Now()

	a.logger.Debug("database get operation started",
		zap.String("username", username),
		zap.String("operation", "get_user"),
	)

	ctx, cancel := context.WithTimeout(ctx, DefaultDBTimeout)
	defer cancel()

	query := `SELECT id, username, email, password FROM users
 			  WHERE username = $1`
	row := a.db.QueryRow(ctx, query, username)
	duration := time.Since(start)

	var user models.User
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Password)
	if err != nil {

		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			a.logger.Error("database query timeout",
				zap.String("email", user.Email),
				zap.String("username", user.Username),
				zap.String("operation", "insert_user"),
				zap.Duration("timeout", DefaultDBTimeout),
				zap.Duration("actual_duration", duration),
				zap.Error(err))
			return models.User{}, fmt.Errorf("database query timeout after %v: %w", DefaultDBTimeout, err)
		}

		if errors.Is(err, pgx.ErrNoRows) {

			a.logger.Error("user not found",
				zap.String("username", username),
				zap.String("operation", "get_user"),
				zap.Duration("duration", duration),
				zap.Error(err),
			)

			return models.User{}, fmt.Errorf("user not found: %w", err)
		}

		a.logger.Error("database select failed",
			zap.String("email", user.Email),
			zap.String("username", user.Username),
			zap.String("operation", "select_user"),
			zap.String("query", "SELECT FROM users"),
			zap.Error(err),
			zap.Duration("db_duration", duration),
		)

		return models.User{}, fmt.Errorf("could not get user: %w", err)
	}

	a.logger.Info("user get successfully",
		zap.Int("user_id", user.ID),
		zap.String("email", user.Email),
		zap.String("username", user.Username),
		zap.String("operation", "select_user"),
		zap.Duration("db_duration", duration),
	)

	if duration > SlowQueryThreshold {
		a.logger.Warn("slow database query detected",
			zap.String("operation", "select_user"),
			zap.Duration("db_duration", duration),
			zap.Duration("threshold", SlowQueryThreshold),
			zap.String("email", user.Email),
		)
	}

	return user, nil
}
