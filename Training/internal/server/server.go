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
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"sync"
	"time"
)

const (
	ShutdownTimeOut = 1 * time.Second
)

type Server struct {
	server       *http.Server
	roomsChecker *room_checker.RoomChecker
}

func NewServer(driver, dsn, appAddress string) Server {
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3333", "http://localhost:3001"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	db, err := db_connection.ConnectDBSQLX(driver, dsn)
	if err != nil {
		logger.Logger.Error(err.Error())
		os.Exit(1)
	}

	trainingRepository := repository.NewTraining(*db)

	trainingUseCase := training_usecase.NewTraining(trainingRepository)

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
	}
}

func (s Server) Run(ctx context.Context, certFile, keyFile string, roomCheckInterval time.Duration) error {

	go func() {
		if err := s.server.ListenAndServe(); err != nil {
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
