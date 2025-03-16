package server

import (
	"Training/internal/background/room_checker"
	"Training/internal/delivery/rest"
	"Training/internal/model"
	"Training/internal/repository"
	"Training/internal/usecase/training_usecase"
	"Training/pkg/db_connection"
	"Training/pkg/logger"
	"context"
	"crypto/tls"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	coachGRPC "github.com/DanKo-code/FitnessCenter-Protobuf/gen/FitnessCenter.protobuf.coach"
	userGRPC "github.com/DanKo-code/FitnessCenter-Protobuf/gen/FitnessCenter.protobuf.user"
)

const (
	ShutdownTimeOut = 1 * time.Second
)

type Server struct {
	server       *http.Server
	roomsChecker *room_checker.RoomChecker
}

func NewServer(driver, dsn, appAddress string) (Server, error) {
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://localhost:3333", "https://localhost:3001"},
		AllowHeaders:     []string{"Content-Type", "Authorization", "Upgrade", "Connection"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		ExposeHeaders:    []string{"Upgrade"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	db, err := db_connection.ConnectDBSQLX(driver, dsn)
	if err != nil {
		logger.Logger.Error(err.Error())
		os.Exit(1)
	}

	trainingRepository := repository.NewTraining(*db)

	connCoach, err := grpc.NewClient(os.Getenv("COACH_SERVICE_PORT"), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Logger.Error("failed to connect to Coach server: %v", err)
		return Server{}, err
	}
	coachClient := coachGRPC.NewCoachClient(connCoach)

	connUer, err := grpc.NewClient(os.Getenv("USER_SERVICE_PORT"), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Logger.Error("failed to connect to User server: %v", err)
		return Server{}, err
	}
	userClient := userGRPC.NewUserClient(connUer)

	trainingUseCase := training_usecase.NewTraining(trainingRepository, coachClient, userClient)

	roomMap := &model.RoomMap{
		Mutex: sync.RWMutex{},
		Map:   make(map[model.RoomMapKey][]model.Participant),
	}
	broadcast := make(chan model.BroadcastMsg)

	rest.RegisterEndpoints(router, trainingUseCase, roomMap, broadcast)

	tlsConfig := &tls.Config{
		InsecureSkipVerify: true, // Игнорировать проверку сертификатов
		ClientAuth:         tls.NoClientCert,
	}

	server := &http.Server{
		Addr:      appAddress,
		Handler:   router,
		TLSConfig: tlsConfig,
	}

	roomsChecker := room_checker.NewRoomChecker(roomMap, trainingUseCase, broadcast)

	return Server{
		server,
		roomsChecker,
	}, nil
}

func (s Server) Run(ctx context.Context, certFile, keyFile string, roomCheckInterval time.Duration) error {

	go func() {
		if err := s.server.ListenAndServeTLS("./certs/cert.crt", "./certs/key.pem"); err != nil {
			logger.Logger.Error(err.Error())
			os.Exit(1)
		}
	}()

	go func() {
		err := s.roomsChecker.Run(ctx, roomCheckInterval)
		if err != nil {
			logger.Logger.Error(err.Error())
			os.Exit(1)
		}
	}()

	go func() {
		s.roomsChecker.RunBroadcaster(ctx)
	}()

	<-ctx.Done()

	ctx, cancel := context.WithTimeout(ctx, ShutdownTimeOut)
	defer cancel()

	return s.server.Shutdown(ctx)
}
