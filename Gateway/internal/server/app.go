package server

import (
	_ "Gateway/docs"
	abonementRest "Gateway/internal/abonement/rest"
	coachRest "Gateway/internal/coach/rest"
	orderRest "Gateway/internal/order/delivery/rest"
	reviewRest "Gateway/internal/review/delivery/rest"
	serviceRest "Gateway/internal/service/delivery/rest"
	ssoRest "Gateway/internal/sso/delivery/rest"
	userRest "Gateway/internal/user/delivery/rest"
	logger "Gateway/pkg/logger"
	"context"
	abonementGRPC "github.com/DanKo-code/FitnessCenter-Protobuf/gen/FitnessCenter.protobuf.abonement"
	coachGRPC "github.com/DanKo-code/FitnessCenter-Protobuf/gen/FitnessCenter.protobuf.coach"
	orderGRPC "github.com/DanKo-code/FitnessCenter-Protobuf/gen/FitnessCenter.protobuf.order"
	reviewGRPC "github.com/DanKo-code/FitnessCenter-Protobuf/gen/FitnessCenter.protobuf.review"
	serviceGRPC "github.com/DanKo-code/FitnessCenter-Protobuf/gen/FitnessCenter.protobuf.service"
	userGRPC "github.com/DanKo-code/FitnessCenter-Protobuf/gen/FitnessCenter.protobuf.user"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net/http"
	"os"
	"os/signal"
	"time"
)

type App struct {
	httpServer      *http.Server
	SSOClient       *grpc.ClientConn
	userClient      *grpc.ClientConn
	abonementClient *abonementGRPC.AbonementClient
	coachClient     *coachGRPC.CoachClient
	reviewClient    *reviewGRPC.ReviewClient
	userClientI     *userGRPC.UserClient
	serviceClient   *serviceGRPC.ServiceClient
	orderClient     *orderGRPC.OrderClient
}

func NewApp() (*App, error) {

	conn, err := grpc.NewClient(os.Getenv("SSO_SERVICE_PORT"), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.ErrorLogger.Printf("failed to connect to SSO server: %v", err)
		return nil, err
	}
	connUer, err := grpc.NewClient(os.Getenv("USER_SERVICE_PORT"), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.ErrorLogger.Printf("failed to connect to User server: %v", err)
		return nil, err
	}
	userClientI := userGRPC.NewUserClient(connUer)

	connAbonement, err := grpc.NewClient(os.Getenv("ABONEMENT_SERVICE_PORT"), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.ErrorLogger.Printf("failed to connect to Abonement server: %v", err)
		return nil, err
	}
	abonementClient := abonementGRPC.NewAbonementClient(connAbonement)

	connCoach, err := grpc.NewClient(os.Getenv("COACH_SERVICE_PORT"), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.ErrorLogger.Printf("failed to connect to Coach server: %v", err)
		return nil, err
	}
	coachClient := coachGRPC.NewCoachClient(connCoach)

	connReview, err := grpc.NewClient(os.Getenv("REVIEW_SERVICE_PORT"), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.ErrorLogger.Printf("failed to connect to Review server: %v", err)
		return nil, err
	}
	reviewClient := reviewGRPC.NewReviewClient(connReview)

	connService, err := grpc.NewClient(os.Getenv("SERVICE_SERVICE_PORT"), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.ErrorLogger.Printf("failed to connect to Service server: %v", err)
		return nil, err
	}
	serviceClient := serviceGRPC.NewServiceClient(connService)

	connOrder, err := grpc.NewClient(os.Getenv("ORDER_SERVICE_PORT"), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.ErrorLogger.Printf("failed to connect to Order server: %v", err)
		return nil, err
	}
	orderClient := orderGRPC.NewOrderClient(connOrder)

	return &App{
		SSOClient:       conn,
		userClient:      connUer,
		abonementClient: &abonementClient,
		coachClient:     &coachClient,
		reviewClient:    &reviewClient,
		userClientI:     &userClientI,
		serviceClient:   &serviceClient,
		orderClient:     &orderClient,
	}, nil
}

func (app *App) Run(port string) error {

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3333", "http://localhost:3001"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	//TODO add cors

	validate := validator.New()

	ssoRest.RegisterHTTPEndpoints(router, validate, app.SSOClient)

	userRest.RegisterHTTPEndpoints(router, validate, app.userClient)

	abonementRest.RegisterHTTPEndpoints(router, app.abonementClient)

	coachRest.RegisterHTTPEndpoints(router, app.coachClient)

	reviewRest.RegisterHTTPEndpoints(router, app.reviewClient, app.userClientI)

	serviceRest.RegisterHTTPEndpoints(router, app.serviceClient)

	orderRest.RegisterHTTPEndpoints(router, app.orderClient)

	router.GET(os.Getenv("SWAGGER_PATH"), ginSwagger.WrapHandler(swaggerFiles.Handler))

	app.httpServer = &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	go func() {
		if err := app.httpServer.ListenAndServe(); err != nil {
			logger.FatalLogger.Printf("Failed to listen and serve: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Interrupt)

	<-quit

	ctx, shutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdown()

	return app.httpServer.Shutdown(ctx)
}
