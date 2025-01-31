package server

import (
	userGRPC "User/internal/delivery/grpc"
	"User/internal/dtos"
	"User/internal/models"
	"User/internal/repository/postgres"
	"User/internal/usecase"
	"User/internal/usecase/localstack_usecase"
	"User/internal/usecase/user_usecase"
	"User/pkg/logger"
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"net"
	"os"
	"os/signal"
	"syscall"
)

type AppGRPC struct {
	gRPCServer   *grpc.Server
	userUseCase  usecase.UserUseCase
	cloudUseCase usecase.CloudUseCase
}

func NewAppGRPC(cloudConfig *models.CloudConfig) (*AppGRPC, error) {

	db := initDB()

	repository := postgres.NewUserRepository(db)

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

	userUseCase := user_usecase.NewUserUseCase(repository, localStackUseCase)

	gRPCServer := grpc.NewServer()

	userGRPC.Register(gRPCServer, userUseCase, localStackUseCase)

	err = insertAdminUser(userUseCase)
	if err != nil {
		return nil, err
	}

	return &AppGRPC{
		gRPCServer:   gRPCServer,
		userUseCase:  userUseCase,
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

func insertAdminUser(userUseCase usecase.UserUseCase) error {

	admins, err := userUseCase.GetAdmins(context.TODO())
	if err != nil {
		return err
	}

	if len(admins) != 0 {
		return nil
	}

	createUserCommand := &dtos.CreateUserCommand{
		ID:       uuid.New(),
		Name:     "Danila",
		Email:    "danilakozlyakovsky@gmail.com",
		Role:     "admin",
		Photo:    "",
		Password: "TankiDanik2003",
	}
	_, err = userUseCase.CreateUser(context.TODO(), createUserCommand)
	if err != nil {
		return err
	}

	logger.InfoLogger.Printf("Admin successfully inserted")
	return nil
}
