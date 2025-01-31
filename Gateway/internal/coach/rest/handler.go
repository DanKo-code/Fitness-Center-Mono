package rest

import (
	"Gateway/internal/coach/dtos"
	"Gateway/pkg/logger"
	"context"
	"fmt"
	coachGRPC "github.com/DanKo-code/FitnessCenter-Protobuf/gen/FitnessCenter.protobuf.coach"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/emptypb"
	"io"
	"net/http"
)

type Handler struct {
	coachClient *coachGRPC.CoachClient
}

func NewHandler(coachClient *coachGRPC.CoachClient) *Handler {
	return &Handler{
		coachClient: coachClient,
	}
}

func (h *Handler) GetCoaches(c *gin.Context) {
	coaches, err := (*h.coachClient).GetCoachesWithServicesWithReviewsWithUsers(c.Request.Context(), &emptypb.Empty{})
	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"err": err,
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"coaches": coaches,
	})
}

func (h *Handler) CreateCoach(c *gin.Context) {

	createCoachCommandAny, exists := c.Get("CreateCoachCommand")
	if !exists {
		logger.ErrorLogger.Printf("Cant find CreateCoachCommand in context")
		return
	}

	createCoachCommand, ok := createCoachCommandAny.(*dtos.CreateCoachCommand)
	if !ok {
		logger.ErrorLogger.Printf("CreateCoachCommand has an invalid type")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "CreateCoachCommand has an invalid type"})
		return
	}

	stream, err := (*h.coachClient).CreateCoach(context.Background())
	if err != nil {
		fmt.Printf("failed to stat file: %v\n", err)
	}

	coachDataForCreate := &coachGRPC.CoachDataForCreate{
		Name:            createCoachCommand.Name,
		Description:     createCoachCommand.Description,
		CoachServiceIds: createCoachCommand.Services,
	}

	createCoachRequestCoachDataForCreate := &coachGRPC.CreateCoachRequest_CoachDataForCreate{
		CoachDataForCreate: coachDataForCreate,
	}

	createCoachRequest := &coachGRPC.CreateCoachRequest{
		Payload: createCoachRequestCoachDataForCreate,
	}

	err = stream.Send(createCoachRequest)
	if err != nil {
		return
	}

	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid form data"})
		return
	}

	photo := form.File["photo"]

	if photo != nil && len(photo) > 0 {
		buffer := make([]byte, 1024*1024)
		file, err := photo[0].Open()
		if err != nil {
			return
		}

		for {
			n, err := file.Read(buffer)
			if err == io.EOF {
				break
			}
			if err != nil {
				return
			}

			err = stream.Send(&coachGRPC.CreateCoachRequest{
				Payload: &coachGRPC.CreateCoachRequest_CoachPhoto{
					CoachPhoto: buffer[:n],
				},
			},
			)
			if err != nil {
				return
			}
		}
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"coach": res.GetCoachWithServices(),
	})
}

func (h *Handler) UpdateCoach(c *gin.Context) {
	updateCoachCommandAny, exists := c.Get("UpdateCoachCommand")
	if !exists {
		logger.ErrorLogger.Printf("Cant find UpdateCoachCommand in context")
		return
	}

	updateCoachCommand, ok := updateCoachCommandAny.(*dtos.UpdateCoachCommand)
	if !ok {
		logger.ErrorLogger.Printf("UpdateCoachCommand has an invalid type")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "UpdateCoachCommand has an invalid type"})
		return
	}

	stream, err := (*h.coachClient).UpdateCoach(context.Background())
	if err != nil {
		fmt.Printf("failed to stat file: %v\n", err)
	}

	coachDataForUpdate := &coachGRPC.CoachDataForUpdate{
		Id:              updateCoachCommand.Id,
		Name:            updateCoachCommand.Name,
		Description:     updateCoachCommand.Description,
		CoachServiceIds: updateCoachCommand.Services,
	}

	updateCoachRequestCoachDataForUpdate := &coachGRPC.UpdateCoachRequest_CoachDataForUpdate{
		CoachDataForUpdate: coachDataForUpdate,
	}

	updateCoachRequest := &coachGRPC.UpdateCoachRequest{
		Payload: updateCoachRequestCoachDataForUpdate,
	}

	err = stream.Send(updateCoachRequest)
	if err != nil {
		logger.ErrorLogger.Printf(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid form data"})
		return
	}

	photo := form.File["photo"]

	if photo != nil && len(photo) > 0 {
		buffer := make([]byte, 1024*1024)
		file, err := photo[0].Open()
		if err != nil {
			logger.ErrorLogger.Printf(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}

		for {
			n, err := file.Read(buffer)
			if err == io.EOF {
				break
			}
			if err != nil {
				logger.ErrorLogger.Printf(err.Error())
				c.JSON(http.StatusInternalServerError, gin.H{"error": err})
				return
			}

			err = stream.Send(&coachGRPC.UpdateCoachRequest{
				Payload: &coachGRPC.UpdateCoachRequest_CoachPhoto{
					CoachPhoto: buffer[:n],
				},
			},
			)
			if err != nil {
				logger.ErrorLogger.Printf(err.Error())
				c.JSON(http.StatusInternalServerError, gin.H{"error": err})
				return
			}
		}
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		logger.ErrorLogger.Printf(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"coach": res.GetCoachWithServices(),
	})
}

func (h *Handler) DeleteCoach(c *gin.Context) {
	id := c.Param("id")

	convertedId, err := uuid.Parse(id)
	if err != nil {

		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": fmt.Errorf("invalid id format")})
		return
	}

	_ = convertedId

	deleteCoachByIdRequest := &coachGRPC.DeleteCoachByIdRequest{
		Id: id,
	}

	deletedCoachRes, err := (*h.coachClient).DeleteCoachById(context.TODO(), deleteCoachByIdRequest)
	if err != nil {
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"coach": deletedCoachRes.GetCoachObject(),
	})
}
