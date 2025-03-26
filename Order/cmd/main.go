package main

import (
	"github.com/DanKo-code/FitnessCenter-Order/internal/server"
	"github.com/DanKo-code/FitnessCenter-Order/pkg/logger"
	"github.com/joho/godotenv"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		logger.FatalLogger.Fatalf("Error loading .env file: %s", err)
	}

	logger.InfoLogger.Printf("Successfully loaded environment variables")

	appGRPC, err := server.NewAppGRPC()
	if err != nil {
		logger.FatalLogger.Fatalf("Error initializing app: %s", err)
	}

	err = appGRPC.Run(os.Getenv("APP_PORT"))
	if err != nil {
		logger.FatalLogger.Fatalf("Error running server")
	}
}
