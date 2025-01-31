package rest

import (
	"Gateway/internal/sso/delivery/rest/middlewares"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc"
)

func RegisterHTTPEndpoints(router *gin.Engine, validator *validator.Validate, ssoClient *grpc.ClientConn) {
	h := NewHandler(ssoClient, validator)

	ssoWithFingerprint := router.Group("/sso")
	{
		ssoWithFingerprint.Use(middlewares.FingerprintMiddleware())
		ssoWithFingerprint.POST("/signUp", h.SignUp)
		ssoWithFingerprint.POST("/signIn", h.SignIn)
		ssoWithFingerprint.POST("/refresh", h.Refresh)
	}

	ssoWithoutFingerprint := router.Group("/sso")
	{
		ssoWithoutFingerprint.POST("/logOut", h.LogOut)
	}
}
