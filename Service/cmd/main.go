package main

import (
	"Service/internal/models"
	"Service/internal/server"
	"Service/pkg/logger"
	"github.com/joho/godotenv"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		logger.FatalLogger.Fatalf("Error loading .env file: %s", err)
	}

	logger.InfoLogger.Printf("Successfully loaded environment variables")

	cloudConfig := &models.CloudConfig{
		EndPoint: os.Getenv("AWS_ENDPOINT"),
		Region:   os.Getenv("AWS_REGION"),
		Bucket:   os.Getenv("AWS_S3_BUCKET"),
		Key:      os.Getenv("AWS_KEY"),
		Secret:   os.Getenv("AWS_SECRET"),
	}

	appGRPC, err := server.NewAppGRPC(cloudConfig)
	if err != nil {
		logger.FatalLogger.Fatalf("Error initializing app: %s", err)
	}

	err = appGRPC.Run(os.Getenv("APP_PORT"))
	if err != nil {
		logger.FatalLogger.Fatalf("Error running server")
	}
}
