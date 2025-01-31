package server

import (
	ssoGRPC "SSO/internal/delivery/grpc"
	"SSO/internal/repository/postgres"
	"SSO/internal/usecase"
	"SSO/internal/usecase/sso_usecase"
	"SSO/pkg/logger"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net"
	"os"
	"os/signal"
)

type AppGRPC struct {
	gRPCServer *grpc.Server
	ssoUseCase usecase.UseCase
	userClient *grpc.ClientConn
}

func NewAppGRPC() (*AppGRPC, error) {

	db := initDB()

	repository := postgres.NewSSORepository(db)

	connUer, err := grpc.NewClient(os.Getenv("USER_SERVICE_PORT"), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.ErrorLogger.Printf("failed to connect to User server: %v", err)
		return nil, err
	}

	useCase := sso_usecase.NewSSOUseCase(repository, connUer)

	return &AppGRPC{
		ssoUseCase: useCase,
		userClient: connUer,
	}, nil
}

func (app *AppGRPC) Run(port string) error {

	app.gRPCServer = grpc.NewServer()

	ssoGRPC.Register(app.gRPCServer, app.ssoUseCase)

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
	signal.Notify(quit, os.Interrupt)

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

	logger.InfoLogger.Printf("Successfully connected to db")

	return db
}
