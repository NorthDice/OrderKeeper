package main

import (
	"OrderKeeper/internal/config"
	"OrderKeeper/internal/handler"
	"OrderKeeper/internal/repository"
	"OrderKeeper/internal/service"
	"OrderKeeper/server"
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	log, err := zap.NewDevelopment()
	if err != nil {
		fmt.Println("unable to initialize logger: %v", err)
		os.Exit(1)
	}
	if err := config.InitConfig(); err != nil {
		fmt.Println("unable to initialize config: %v", err)
		log.Fatal("error initializing config")
		os.Exit(1)
	}
	if err := godotenv.Load(); err != nil {
		log.Fatal("error loading .env file")
		os.Exit(1)
	}

	db, err := repository.NewPostgresDB(context.Background(), repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
	})
	if err != nil {
		fmt.Println("error initializing db connection: %v", err)
		log.Fatal("error initializing postgres db")
		os.Exit(1)
	}
	repo := repository.NewRepository(db, log)
	services := service.NewService(repo, log)
	handlers := handler.NewHandler(services, log)

	srv := new(server.Server)
	go func() {
		if err := srv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
			log.Fatal("error running server", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	log.Info("Keeper shutting down")

	if err := srv.Shutdown(context.Background()); err != nil {
		log.Error("err occured while shutting down: %s")
	}
	db.Close()
}
