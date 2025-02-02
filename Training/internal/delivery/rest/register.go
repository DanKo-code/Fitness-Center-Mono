package rest

import (
	"Training/internal/common_middlewares/middlewares"
	"github.com/gin-gonic/gin"
)

func RegisterEndpoints(router *gin.Engine, useCase trainingUseCase) {
	h := NewHandler(useCase)

	authorized := router.Group("/", middlewares.VerifyAccessTokenMiddleware())
	authorized.POST("/training/book", h.Insert)
}
