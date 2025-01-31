package rest

import (
	"Gateway/internal/common_middlewares/middlewares"
	orderGRPC "github.com/DanKo-code/FitnessCenter-Protobuf/gen/FitnessCenter.protobuf.order"
	"github.com/gin-gonic/gin"
)

func RegisterHTTPEndpoints(router *gin.Engine, orderClient *orderGRPC.OrderClient) {
	h := NewHandler(orderClient)

	router.POST("/checkout-session-completed", h.HandleCheckoutSessionCompleted)
	router.POST("/create-checkout-session", middlewares.VerifyAccessTokenMiddleware(), middlewares.IsClientMiddleware(), h.CreateCheckoutSession)
	router.GET("/orders/:userId", middlewares.VerifyAccessTokenMiddleware(), h.GetUserOrders)
}
