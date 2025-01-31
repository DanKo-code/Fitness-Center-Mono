package rest

import (
	"Gateway/internal/common_middlewares/middlewares"
	reviewGRPC "github.com/DanKo-code/FitnessCenter-Protobuf/gen/FitnessCenter.protobuf.review"
	userGRPC "github.com/DanKo-code/FitnessCenter-Protobuf/gen/FitnessCenter.protobuf.user"
	"github.com/gin-gonic/gin"
)

func RegisterHTTPEndpoints(
	router *gin.Engine,
	reviewClient *reviewGRPC.ReviewClient,
	userClient *userGRPC.UserClient,
) {
	h := NewHandler(reviewClient, userClient)

	router.POST("/reviews", middlewares.VerifyAccessTokenMiddleware(), middlewares.IsClientMiddleware(), h.CreateCoachReview)
}
