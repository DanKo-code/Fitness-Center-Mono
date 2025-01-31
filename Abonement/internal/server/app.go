package server

import (
	"context"
	"fmt"
	userGRPC "github.com/DanKo-code/Fitness-Center-Abonement/internal/delivery/grpc"
	"github.com/DanKo-code/Fitness-Center-Abonement/internal/models"
	"github.com/DanKo-code/Fitness-Center-Abonement/internal/repository/postgres"
	"github.com/DanKo-code/Fitness-Center-Abonement/internal/usecase"
	user_usecase "github.com/DanKo-code/Fitness-Center-Abonement/internal/usecase/abonement_usecase"
	"github.com/DanKo-code/Fitness-Center-Abonement/internal/usecase/localstack_usecase"
	"github.com/DanKo-code/Fitness-Center-Abonement/internal/usecase/stripe_usecase"
	"github.com/DanKo-code/Fitness-Center-Abonement/pkg/logger"
	serviceGRPC "github.com/DanKo-code/FitnessCenter-Protobuf/gen/FitnessCenter.protobuf.service"
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
	gRPCServer       *grpc.Server
	abonementUseCase usecase.AbonementUseCase
	cloudUseCase     usecase.CloudUseCase
}

func NewAppGRPC(cloudConfig *models.CloudConfig) (*AppGRPC, error) {

	db := initDB()

	repository := postgres.NewAbonementRepository(db)

	connService, err := grpc.NewClient(os.Getenv("SERVICE_SERVICE_PORT"), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.ErrorLogger.Printf("failed to connect to Service server: %v", err)
		return nil, err
	}

	serviceClient := serviceGRPC.NewServiceClient(connService)

	suc := stripe_usecase.NewStripeUseCase(os.Getenv("STRIPE_KEY"))

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

	abonementUseCase := user_usecase.NewAbonementUseCase(repository, &serviceClient, suc, localStackUseCase)

	gRPCServer := grpc.NewServer()

	userGRPC.RegisterAbonementServer(gRPCServer, abonementUseCase, localStackUseCase, &serviceClient)

	return &AppGRPC{
		gRPCServer:       gRPCServer,
		abonementUseCase: abonementUseCase,
		cloudUseCase:     localStackUseCase,
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
