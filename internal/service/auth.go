package service

import (
	"OrderKeeper/internal/models"
	"OrderKeeper/internal/repository"
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"os"
	"time"
)

const (
	tokenTTL = 12 * time.Hour
)

type tokenClaims struct {
	jwt.RegisteredClaims
	UserID int `json:"user_id"`
}
type AuthorizationService struct {
	repo   repository.Authorization
	logger *zap.Logger
}

func NewAuthorizationService(repository repository.Authorization, logger *zap.Logger) *AuthorizationService {
	return &AuthorizationService{
		repo:   repository,
		logger: logger,
	}
}

func (a *AuthorizationService) CreateUser(ctx context.Context, user models.User) (int, error) {
	start := time.Now()

	a.logger.Info("user creation process started",
		zap.String("email", user.Email),
		zap.String("username", user.Username),
	)

	user.Password = generatePasswordHash(user.Password)
	id, err := a.repo.CreateUser(ctx, user)
	if err != nil {
		a.logger.Error("failed to create user", zap.Error(err),
			zap.String("email", user.Email),
			zap.String("username", user.Username),
			zap.Error(err),
			zap.Duration("total_duration", time.Since(start)),
		)
		return 0, fmt.Errorf("failed to create user: %w", err)
	}

	a.logger.Info("user created successfully",
		zap.Int("user_id", id),
		zap.String("email", user.Email),
		zap.String("username", user.Username),
		zap.Duration("total_service_duration", time.Since(start)),
	)

	return id, nil
}

func (a *AuthorizationService) GenerateToken(ctx context.Context, username, password string) (string, error) {
	start := time.Now()
	a.logger.Info("user get process started",
		zap.String("username", username),
	)

	user, err := a.repo.GetUser(ctx, username, password)
	if err != nil {
		a.logger.Error("failed to get user",
			zap.String("username", username),
			zap.Error(err),
			zap.Duration("total_duration", time.Since(start)),
		)
		return "", fmt.Errorf("failed to get user: %w", err)
	}

	a.logger.Info("token generation process started",
		zap.String("username", username),
		zap.Int("user_id", user.ID),
	)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		UserID: user.ID,
	})

	a.logger.Info("token generated successfully",
		zap.String("username", username),
		zap.Int("user_id", user.ID),
		zap.Duration("total_service_duration", time.Since(start)),
	)

	return token.SignedString([]byte(os.Getenv("SIGNING_KEY")))
}
func (a *AuthorizationService) ParseToken(ctx context.Context, token string) (int, error) {
	parsedToken, err := jwt.ParseWithClaims(token, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("SIGNING_KEY")), nil
	})
	if err != nil {
		return 0, err
	}

	claims, ok := parsedToken.Claims.(*tokenClaims)
	if !ok || !parsedToken.Valid {
		return 0, fmt.Errorf("invalid token")
	}

	return claims.UserID, nil
}

func generatePasswordHash(password string) string {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		fmt.Errorf("could not generate password: %v", err)
	}

	return string(hashedPassword)
}
