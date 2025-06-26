package main

import (
	"OrderKeeper/internal/config"
	"OrderKeeper/internal/handler"
	"OrderKeeper/internal/repository/cache"
	"OrderKeeper/internal/repository/postgres"
	"OrderKeeper/internal/service"
	"OrderKeeper/server"
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	logger := initLogger()

	if err := config.InitConfig(); err != nil {
		fmt.Println("unable to initialize config: %v", err)
		logger.Fatal("error initializing config")
		os.Exit(1)
	}

	if err := godotenv.Load(); err != nil {
		logger.Fatal("error loading .env file")
		os.Exit(1)
	}

	db, err := postgres.NewPostgresDB(context.Background(), postgres.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
	})
	if err != nil {
		logger.Fatal("error initializing postgres db")
		os.Exit(1)
	}
	defer db.Close()
	logger.Info("Postgres DB initialized successfully")

	var repo *postgres.Repository

	redisEnabled := viper.GetBool("redis.enable")
	if redisEnabled {
		redisCache, err := cache.NewRedisCache(context.Background(), cache.RedisConfig{
			Address:  viper.GetString("redis.address"),
			Password: viper.GetString("redis.password"),
			Database: viper.GetInt("redis.database"),
		}, logger)
		if err != nil {
			logger.Error("error initializing redis, falling back to non-cached repository", zap.Error(err))
			repo = postgres.NewRepository(db, logger)
		} else {
			logger.Info("Redis initialized successfully, using cached repository")
			repo = postgres.NewCachedRepository(db, &redisCache, logger)
		}
	} else {
		logger.Info("Redis disabled, using non-cached repository")
		repo = postgres.NewRepository(db, logger)
	}

	services := service.NewService(repo, logger)
	handlers := handler.NewHandler(services, logger)

	srv := new(server.Server)
	go func() {
		if err := srv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
			logger.Fatal("error running server", zap.Error(err))
		}
	}()

	logger.Info("Server started on port", zap.String("port", viper.GetString("port")))

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Keeper shutting down")

	if err := srv.Shutdown(context.Background()); err != nil {
		logger.Error("error occurred while shutting down", zap.Error(err))
	}

	logger.Info("Server exited")
}

func initLogger() *zap.Logger {
	cfg := zap.NewProductionConfig()

	cfg.EncoderConfig.TimeKey = "time"
	cfg.EncoderConfig.LevelKey = "level"
	cfg.EncoderConfig.MessageKey = "msg"
	cfg.EncoderConfig.CallerKey = "caller"
	cfg.EncoderConfig.StacktraceKey = "stacktrace"
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	cfg.EncoderConfig.EncodeLevel = zapcore.LowercaseLevelEncoder

	cfg.InitialFields = map[string]interface{}{
		"service": "myapp",
		"version": "1.0.0",
	}

	logger, err := cfg.Build()
	if err != nil {
		log.Fatal("Failed to initialize logger:", err)
	}

	return logger
}
