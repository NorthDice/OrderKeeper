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
		logger.Warn("no .env file found, using environment variables")
	}

	db, err := postgres.NewPostgresDB(context.Background(), postgres.Config{
		Host:     getConfigString("db.host", "DB_HOST"),
		Port:     getConfigString("db.port", "DB_PORT"),
		Username: getConfigString("db.username", "DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   getConfigString("db.dbname", "DB_NAME"),
		SSLMode:  getConfigString("db.sslmode", "DB_SSLMODE"),
	})
	if err != nil {
		logger.Fatal("error initializing postgres db", zap.Error(err))
		os.Exit(1)
	}
	defer db.Close()
	logger.Info("Postgres DB initialized successfully")

	var repo *postgres.Repository

	redisEnabled := getConfigBool("redis.enable", "REDIS_ENABLE")
	if redisEnabled {
		// Try environment variables first, then config file
		redisAddr := os.Getenv("REDIS_ADDR")
		if redisAddr == "" {
			redisHost := getConfigString("redis.host", "REDIS_HOST")
			redisPort := getConfigString("redis.port", "REDIS_PORT")
			redisAddr = fmt.Sprintf("%s:%s", redisHost, redisPort)
		}

		redisPassword := os.Getenv("REDIS_PASSWORD")
		redisDB := getConfigInt("redis.db", "REDIS_DB")

		logger.Info("Attempting to connect to Redis", zap.String("address", redisAddr))

		redisCache, err := cache.NewRedisCache(context.Background(), cache.RedisConfig{
			Address:  redisAddr,
			Password: redisPassword,
			Database: redisDB,
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
		port := getConfigString("port", "PORT")
		if port == "" {
			port = "8080"
		}
		if err := srv.Run(":"+port, handlers.InitRoutes()); err != nil {
			logger.Fatal("error running server", zap.Error(err))
		}
	}()

	port := getConfigString("port", "PORT")
	if port == "" {
		port = "8080"
	}
	logger.Info("Server started on port", zap.String("port", port))

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Keeper shutting down")

	if err := srv.Shutdown(context.Background()); err != nil {
		logger.Error("error occurred while shutting down", zap.Error(err))
	}

	logger.Info("Server exited")
}

// Helper functions to prioritize environment variables over config file
func getConfigString(configKey, envKey string) string {
	if envVal := os.Getenv(envKey); envVal != "" {
		return envVal
	}
	return viper.GetString(configKey)
}

func getConfigBool(configKey, envKey string) bool {
	if envVal := os.Getenv(envKey); envVal != "" {
		return envVal == "true" || envVal == "1"
	}
	return viper.GetBool(configKey)
}

func getConfigInt(configKey, envKey string) int {
	if envVal := os.Getenv(envKey); envVal != "" {
		// You might want to add proper string to int conversion with error handling
		if envVal == "0" {
			return 0
		}
		return 1 // Default for non-zero values
	}
	return viper.GetInt(configKey)
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
