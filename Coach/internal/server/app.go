package server

import (
	"context"
	"fmt"
	coachGRPC "github.com/DanKo-code/FitnessCenter-Coach/internal/delivery/grpc"
	"github.com/DanKo-code/FitnessCenter-Coach/internal/models"
	"github.com/DanKo-code/FitnessCenter-Coach/internal/repository/postgres"
	"github.com/DanKo-code/FitnessCenter-Coach/internal/usecase"
	"github.com/DanKo-code/FitnessCenter-Coach/internal/usecase/coach_usecase"
	"github.com/DanKo-code/FitnessCenter-Coach/internal/usecase/localstack_usecase"
	"github.com/DanKo-code/FitnessCenter-Coach/pkg/logger"
	reviewGRPC "github.com/DanKo-code/FitnessCenter-Protobuf/gen/FitnessCenter.protobuf.review"
	serviceGRPC "github.com/DanKo-code/FitnessCenter-Protobuf/gen/FitnessCenter.protobuf.service"
	userGRPC "github.com/DanKo-code/FitnessCenter-Protobuf/gen/FitnessCenter.protobuf.user"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net"
	"os"
	"os/signal"
	"syscall"
)

type AppGRPC struct {
	gRPCServer   *grpc.Server
	coachUseCase usecase.CoachUseCase
	cloudUseCase usecase.CloudUseCase
}

func NewAppGRPC(cloudConfig *models.CloudConfig) (*AppGRPC, error) {

	db := initDB()

	repository := postgres.NewCoachRepository(db)

	connService, err := grpc.NewClient(os.Getenv("SERVICE_SERVICE_PORT"), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.ErrorLogger.Printf("failed to connect to Service server: %v", err)
		return nil, err
	}
	serviceClient := serviceGRPC.NewServiceClient(connService)

	connReview, err := grpc.NewClient(os.Getenv("REVIEW_SERVICE_PORT"), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.ErrorLogger.Printf("failed to connect to Review server: %v", err)
		return nil, err
	}
	reviewClient := reviewGRPC.NewReviewClient(connReview)

	connUser, err := grpc.NewClient(os.Getenv("USER_SERVICE_PORT"), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.ErrorLogger.Printf("failed to connect to User server: %v", err)
		return nil, err
	}
	userClient := userGRPC.NewUserClient(connUser)

	awsCfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(cloudConfig.Region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(cloudConfig.Key, cloudConfig.Secret, "")),
	)
	if err != nil {
		logger.FatalLogger.Fatalf("failed loading config, %v", err)
		return nil, err
	}

	client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		o.UsePathStyle = true
		o.BaseEndpoint = aws.String(cloudConfig.EndPoint)
	})

	localStackUseCase := localstack_usecase.NewLocalstackUseCase(client, cloudConfig)

	coachUseCase := user_usecase.NewCoachUseCase(repository, &serviceClient, &reviewClient, &userClient, localStackUseCase)

	gRPCServer := grpc.NewServer()

	coachGRPC.RegisterCoachServer(gRPCServer, coachUseCase, localStackUseCase, &serviceClient)

	return &AppGRPC{
		gRPCServer:   gRPCServer,
		coachUseCase: coachUseCase,
		cloudUseCase: localStackUseCase,
	}, nil
}

func (app *AppGRPC) Run(port string) error {

	listen, err := net.Listen(os.Getenv("APP_GRPC_PROTOCOL"), port)
	if err != nil {
		logger.ErrorLogger.Printf("Failed to listen: %v", err)
		return err
	}

	logger.InfoLogger.Printf("Starting gRPC server on port %s", port)

	go func() {
		if err = app.gRPCServer.Serve(listen); err != nil {
			logger.FatalLogger.Fatalf("Failed to serve: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	logger.InfoLogger.Printf("stopping gRPC server %s", port)
	app.gRPCServer.GracefulStop()

	return nil
}

func initDB() *sqlx.DB {

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_SLLMODE"),
	)

	db, err := sqlx.Connect(os.Getenv("DB_DRIVER"), dsn)
	if err != nil {
		logger.FatalLogger.Fatalf("Database connection failed: %s", err)
	}

	logger.InfoLogger.Println("Successfully connected to db")

	return db
}
