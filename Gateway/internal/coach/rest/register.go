package rest

import (
	coachMW "Gateway/internal/coach/middlewares"
	"Gateway/internal/common_middlewares/middlewares"
	coachGRPC "github.com/DanKo-code/FitnessCenter-Protobuf/gen/FitnessCenter.protobuf.coach"
	"github.com/gin-gonic/gin"
)

func RegisterHTTPEndpoints(router *gin.Engine, coachClient *coachGRPC.CoachClient) {
	h := NewHandler(coachClient)

	authorized := router.Group("/", middlewares.VerifyAccessTokenMiddleware())

	authorized.GET("/coaches", h.GetCoaches)
	authorized.POST("/coaches", middlewares.IsAdminMiddleware(), coachMW.ValidateCreateCoachMW(), h.CreateCoach)
	authorized.PUT("/coaches", middlewares.IsAdminMiddleware(), coachMW.ValidateUpdateCoachMW(), h.UpdateCoach)
	authorized.DELETE("/coaches/:id", middlewares.IsAdminMiddleware(), h.DeleteCoach)
}
