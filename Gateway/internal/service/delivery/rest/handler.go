package rest

import (
	"context"
	serviceGRPC "github.com/DanKo-code/FitnessCenter-Protobuf/gen/FitnessCenter.protobuf.service"
	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/types/known/emptypb"
	"net/http"
)

type Handler struct {
	serviceClient *serviceGRPC.ServiceClient
}

func NewHandler(serviceClient *serviceGRPC.ServiceClient) *Handler {
	return &Handler{
		serviceClient: serviceClient,
	}
}

func (h *Handler) GetServices(c *gin.Context) {
	services, err := (*h.serviceClient).GetServices(context.TODO(), &emptypb.Empty{})
	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"err": err,
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"services": services,
	})
}
