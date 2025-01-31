package server

import (
	"context"
	"fmt"
	"github.com/DanKo-code/FitnessCenter-Order/internal/background"
	"github.com/DanKo-code/FitnessCenter-Order/internal/background/order_background"
	orderGRPC "github.com/DanKo-code/FitnessCenter-Order/internal/delivery/grpc"
	"github.com/DanKo-code/FitnessCenter-Order/internal/repository/postgres"
	"github.com/DanKo-code/FitnessCenter-Order/internal/usecase/order_usecase"
	"github.com/DanKo-code/FitnessCenter-Order/pkg/logger"
	abonementGRPC "github.com/DanKo-code/FitnessCenter-Protobuf/gen/FitnessCenter.protobuf.abonement"
	serviceGRPC "github.com/DanKo-code/FitnessCenter-Protobuf/gen/FitnessCenter.protobuf.service"
	userGRPC "github.com/DanKo-code/FitnessCenter-Protobuf/gen/FitnessCenter.protobuf.user"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	interval = 1 * time.Minute / 3
)

type AppGRPC struct {
	gRPCServer *grpc.Server
	oc         background.OrderExpiredChecker
}

func NewAppGRPC() (*AppGRPC, error) {

	db := initDB()

	connAbonement, err := grpc.NewClient(os.Getenv("ABONEMENT_SERVICE_PORT"), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.ErrorLogger.Printf("failed to connect to abonement server: %v", err)
		return nil, err
	}
	connUser, err := grpc.NewClient(os.Getenv("USER_SERVICE_PORT"), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.ErrorLogger.Printf("failed to connect to user server: %v", err)
		return nil, err
	}

	connService, err := grpc.NewClient(os.Getenv("SERVICE_SERVICE_PORT"), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.ErrorLogger.Printf("failed to connect to service server: %v", err)
		return nil, err
	}

	abonementClient := abonementGRPC.NewAbonementClient(connAbonement)
	userClient := userGRPC.NewUserClient(connUser)
	serviceClient := serviceGRPC.NewServiceClient(connService)

	repository := postgres.NewOrderRepository(db)

	OrderUseCase := order_usecase.NewOrderUseCase(repository, &abonementClient, &userClient, &serviceClient)

	gRPCServer := grpc.NewServer()

	orderGRPC.Register(gRPCServer, OrderUseCase)

	oc := order_background.NewOrderExpiredChecker(OrderUseCase)

	return &AppGRPC{
		gRPCServer: gRPCServer,
		oc:         oc,
	}, nil
}

func (app *AppGRPC) Run(port string) error {

	listen, err := net.Listen(os.Getenv("APP_GRPC_PROTOCOL"), port)
	if err != nil {
		logger.ErrorLogger.Printf("Failed to listen: %v", err)
		return err
	}

	logger.InfoLogger.Printf("Starting gRPC server on port %s", port)

	stopChecker := make(chan struct{})

	go func() {
		if err = app.gRPCServer.Serve(listen); err != nil {
			logger.FatalLogger.Fatalf("Failed to serve: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	go app.oc.StartOrderExpiredChecker(context.TODO(), interval, stopChecker)

	<-quit

	close(stopChecker)

	logger.InfoLogger.Printf("stopping gRPC server %s", port)

	go app.gRPCServer.GracefulStop()

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
