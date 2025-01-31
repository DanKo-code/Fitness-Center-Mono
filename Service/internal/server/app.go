package server

import (
	serviceGRPC "Service/internal/delivery/grpc"
	"Service/internal/dtos"
	"Service/internal/models"
	"Service/internal/repository/postgres"
	"Service/internal/usecase"
	"Service/internal/usecase/localstack_usecase"
	"Service/internal/usecase/service_usecase"
	"Service/pkg/logger"
	"context"
	"fmt"
	abonementGRPC "github.com/DanKo-code/FitnessCenter-Protobuf/gen/FitnessCenter.protobuf.abonement"
	coachGRPC "github.com/DanKo-code/FitnessCenter-Protobuf/gen/FitnessCenter.protobuf.coach"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io/ioutil"
	"net"
	"os"
	"os/signal"
	"syscall"
)

type AppGRPC struct {
	gRPCServer     *grpc.Server
	serviceUseCase usecase.ServiceUseCase
	cloudUseCase   usecase.CloudUseCase
	coachClient    *coachGRPC.CoachClient
}

func NewAppGRPC(cloudConfig *models.CloudConfig) (*AppGRPC, error) {

	db := initDB()

	connCoach, err := grpc.NewClient(os.Getenv("COACH_SERVICE_PORT"), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.ErrorLogger.Printf("failed to connect to coach server: %v", err)
		return nil, err
	}

	connAbonement, err := grpc.NewClient(os.Getenv("ABONEMENT_SERVICE_PORT"), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.ErrorLogger.Printf("failed to connect to abonement server: %v", err)
		return nil, err
	}

	coachClient := coachGRPC.NewCoachClient(connCoach)
	abonementClient := abonementGRPC.NewAbonementClient(connAbonement)

	repository := postgres.NewServiceRepository(db)

	serviceUseCase := service_usecase.NewServiceUseCase(repository, &coachClient, &abonementClient)

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

	gRPCServer := grpc.NewServer()

	serviceGRPC.Register(gRPCServer, serviceUseCase, localStackUseCase)

	//to do initial insert if no data
	err = insertInitServices(serviceUseCase, localStackUseCase)
	if err != nil {
		return nil, err
	}

	return &AppGRPC{
		gRPCServer:     gRPCServer,
		serviceUseCase: serviceUseCase,
		cloudUseCase:   localStackUseCase,
		coachClient:    &coachClient,
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

func insertInitServices(serviceUseCase usecase.ServiceUseCase, cloudUseCase usecase.CloudUseCase) error {

	services, err := serviceUseCase.GetServices(context.TODO())
	if err != nil {
		return err
	}

	if len(services) != 0 {
		return nil
	}

	//gym insert
	gymPhotoBytes, err := readImageToBytes("internal/images/gym.png")
	if err != nil {
		return err
	}
	randomGymID := uuid.New().String()
	gymUrl, err := cloudUseCase.PutObject(context.TODO(), gymPhotoBytes, "service/"+randomGymID)
	if err != nil {
		return err
	}
	gymServiceReq := &dtos.CreateServiceCommand{
		Id:    uuid.New(),
		Title: "gym",
		Photo: gymUrl,
	}
	_, err = serviceUseCase.CreateService(context.TODO(), gymServiceReq)
	if err != nil {
		return err
	}

	//gym sauna
	saunaPhotoBytes, err := readImageToBytes("internal/images/sauna.png")
	if err != nil {
		return err
	}
	randomSaunaID := uuid.New().String()
	saunaUrl, err := cloudUseCase.PutObject(context.TODO(), saunaPhotoBytes, "service/"+randomSaunaID)
	if err != nil {
		return err
	}
	saunaServiceReq := &dtos.CreateServiceCommand{
		Id:    uuid.New(),
		Title: "sauna",
		Photo: saunaUrl,
	}
	_, err = serviceUseCase.CreateService(context.TODO(), saunaServiceReq)
	if err != nil {
		return err
	}

	//gym swimming-pool
	swimmingPoolPhotoBytes, err := readImageToBytes("internal/images/swimming-pool.png")
	if err != nil {
		return err
	}
	randomSwimmingPoolID := uuid.New().String()
	swimmingPoolUrl, err := cloudUseCase.PutObject(context.TODO(), swimmingPoolPhotoBytes, "service/"+randomSwimmingPoolID)
	if err != nil {
		return err
	}
	swimmingPoolServiceReq := &dtos.CreateServiceCommand{
		Id:    uuid.New(),
		Title: "swimming-pool",
		Photo: swimmingPoolUrl,
	}
	_, err = serviceUseCase.CreateService(context.TODO(), swimmingPoolServiceReq)
	if err != nil {
		return err
	}

	logger.InfoLogger.Printf("Init Services successfully inserted")
	return nil
}

func readImageToBytes(filePath string) ([]byte, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("не удалось открыть файл: %w", err)
	}
	defer file.Close()

	imageBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("не удалось прочитать файл: %w", err)
	}

	return imageBytes, nil
}
