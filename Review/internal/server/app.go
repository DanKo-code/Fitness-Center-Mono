package server

import (
	"fmt"
	coachGRPC "github.com/DanKo-code/FitnessCenter-Protobuf/gen/FitnessCenter.protobuf.coach"
	userGRPC "github.com/DanKo-code/FitnessCenter-Protobuf/gen/FitnessCenter.protobuf.user"
	reviewGRPC "github.com/DanKo-code/FitnessCenter-Review/internal/delivery/grpc"
	"github.com/DanKo-code/FitnessCenter-Review/internal/repository/postgres"
	"github.com/DanKo-code/FitnessCenter-Review/internal/usecase/review_usecase"
	"github.com/DanKo-code/FitnessCenter-Review/pkg/logger"
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
	gRPCServer  *grpc.Server
	coachClient *coachGRPC.CoachClient
}

func NewAppGRPC() (*AppGRPC, error) {

	db := initDB()

	connCoach, err := grpc.NewClient(os.Getenv("COACH_SERVICE_PORT"), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.ErrorLogger.Printf("failed to connect to coach server: %v", err)
		return nil, err
	}
	connUser, err := grpc.NewClient(os.Getenv("USER_SERVICE_PORT"), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.ErrorLogger.Printf("failed to connect to user server: %v", err)
		return nil, err
	}

	coachClient := coachGRPC.NewCoachClient(connCoach)
	userClient := userGRPC.NewUserClient(connUser)

	repository := postgres.NewReviewRepository(db)

	reviewUseCase := review_usecase.NewReviewUseCase(repository, &coachClient, &userClient)

	gRPCServer := grpc.NewServer()

	reviewGRPC.Register(gRPCServer, reviewUseCase)

	return &AppGRPC{
		gRPCServer: gRPCServer,
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
