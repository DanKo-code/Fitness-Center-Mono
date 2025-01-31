package main

import (
	"SSO/internal/server"
	"SSO/pkg/logger"
	"fmt"
	"github.com/joho/godotenv"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		logger.FatalLogger.Fatalf(fmt.Sprintf("Error loading .env file: %s", err))
	}

	logger.InfoLogger.Printf("Successfully loaded environment variables")

	appGRPC, err := server.NewAppGRPC()
	if err != nil {
		logger.FatalLogger.Fatalf(fmt.Sprintf("Error initializing app grpc: %s", err))
	}

	err = appGRPC.Run(os.Getenv("APP_PORT"))
	if err != nil {
		logger.FatalLogger.Printf("Error running server")
	}
}
