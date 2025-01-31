package rest

import (
	"Gateway/internal/common_middlewares/middlewares"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc"
)

func RegisterHTTPEndpoints(router *gin.Engine, validator *validator.Validate, ssoClient *grpc.ClientConn) {
	h := NewHandler(ssoClient, validator)

	authorized := router.Group("/", middlewares.VerifyAccessTokenMiddleware())

	authorized.PUT("/users/:id", middlewares.IsClientMiddleware(), h.UpdateUser)
	authorized.GET("/users", middlewares.IsAdminMiddleware(), h.GetClients)
	authorized.DELETE("/users/:id", middlewares.IsAdminMiddleware(), h.DeleteClientById)
}
