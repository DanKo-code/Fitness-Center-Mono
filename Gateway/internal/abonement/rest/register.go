package rest

import (
	abonementMW "Gateway/internal/abonement/middlewares"
	"Gateway/internal/common_middlewares/middlewares"
	abonementGRPC "github.com/DanKo-code/FitnessCenter-Protobuf/gen/FitnessCenter.protobuf.abonement"
	"github.com/gin-gonic/gin"
)

func RegisterHTTPEndpoints(router *gin.Engine, abonementClient *abonementGRPC.AbonementClient) {
	h := NewHandler(abonementClient)

	authorized := router.Group("/", middlewares.VerifyAccessTokenMiddleware())

	authorized.GET("/abonements", h.GetAbonements)
	authorized.POST("/abonements", middlewares.IsAdminMiddleware(), abonementMW.ValidateCreateAbonementMW(), h.CreateAbonement)
	authorized.PUT("/abonements", middlewares.IsAdminMiddleware(), abonementMW.ValidateUpdateAbonementMW(), h.UpdateAbonement)
	authorized.DELETE("/abonements/:id", middlewares.IsAdminMiddleware(), h.DeleteAbonement)
}
