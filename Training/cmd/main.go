package main

import (
	"Training/internal/server"
	"Training/pkg/logger"
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	roomCheckerInterval = 5 * time.Second
)

func main() {

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	err := godotenv.Load()
	if err != nil {
		logger.Logger.Error(err.Error())
		os.Exit(1)
	}

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_SLLMODE"),
	)

	newServer := server.NewServer(os.Getenv("DB_DRIVER"), dsn, os.Getenv("APP_ADDRESS"))

	err = newServer.Run(ctx, os.Getenv("APP_CERT_FILE"), os.Getenv("APP_KEY_FILE"), roomCheckerInterval)
	if err != nil {
		logger.Logger.Error(err.Error())
		os.Exit(1)
	}
}
