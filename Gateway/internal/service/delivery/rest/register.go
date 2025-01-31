package rest

import (
	serviceGRPC "github.com/DanKo-code/FitnessCenter-Protobuf/gen/FitnessCenter.protobuf.service"
	"github.com/gin-gonic/gin"
)

func RegisterHTTPEndpoints(
	router *gin.Engine,
	serviceClient *serviceGRPC.ServiceClient,
) {
	h := NewHandler(serviceClient)

	router.GET("/services", h.GetServices)
}
