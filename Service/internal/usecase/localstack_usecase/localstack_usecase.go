package localstack_usecase

import (
	"Service/internal/models"
	"Service/pkg/logger"
	"bytes"
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"strings"
)

type LocalstackUseCase struct {
	client *s3.Client
	config *models.CloudConfig
}

func NewLocalstackUseCase(client *s3.Client, config *models.CloudConfig) *LocalstackUseCase {
	return &LocalstackUseCase{
		client: client,
		config: config,
	}
}

func (luc *LocalstackUseCase) PutObject(ctx context.Context, object []byte, name string) (string, error) {
	_, err := luc.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(luc.config.Bucket),
		Key:    aws.String(name),
		Body:   bytes.NewReader(object),
	})
	if err != nil {
		logger.ErrorLogger.Printf("Failed to put object: %v", err)
		return "", err
	}

	fileURL := fmt.Sprintf("%s/%s/%s", luc.config.EndPoint, luc.config.Bucket, name)

	//change to localhost
	fileURL = strings.Replace(fileURL, "localstack", "localhost", 1)

	return fileURL, nil
}

func (luc *LocalstackUseCase) DeleteObject(ctx context.Context, name string) error {
	_, err := luc.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(luc.config.Bucket),
		Key:    aws.String(name),
	})
	if err != nil {
		logger.ErrorLogger.Printf("Failed to delete object: %v", err)
		return err
	}

	return nil
}
