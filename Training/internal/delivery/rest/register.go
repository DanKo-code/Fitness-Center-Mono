package rest

import (
	"Training/internal/common_middlewares/middlewares"
	"Training/internal/model"
	"github.com/gin-gonic/gin"
)

func RegisterEndpoints(router *gin.Engine, useCase trainingUseCase, roomMap *model.RoomMap, broadcast chan model.BroadcastMsg) {
	h := NewHandler(useCase, roomMap, broadcast)

	authorized := router.Group("/", middlewares.VerifyAccessTokenMiddleware())
	authorized.POST("/training/book", h.Insert)
	router.GET("/training/join/:roomId", h.Join)
	router.GET("/training/:day", h.GetTrainingsByDay)
}
