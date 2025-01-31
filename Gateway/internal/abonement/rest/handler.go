package rest

import (
	"Gateway/internal/abonement/dtos"
	"Gateway/pkg/logger"
	"context"
	"fmt"
	abonementGRPC "github.com/DanKo-code/FitnessCenter-Protobuf/gen/FitnessCenter.protobuf.abonement"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/emptypb"
	"io"
	"net/http"
)

type Handler struct {
	abonementClient *abonementGRPC.AbonementClient
}

func NewHandler(abonementClient *abonementGRPC.AbonementClient) *Handler {
	return &Handler{
		abonementClient: abonementClient,
	}
}

func (h *Handler) GetAbonements(c *gin.Context) {
	abonements, err := (*h.abonementClient).GetAbonementsWithServices(c.Request.Context(), &emptypb.Empty{})
	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"err": err,
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"abonements": abonements,
	})
}

func (h *Handler) CreateAbonement(c *gin.Context) {

	createAbonementCommandAny, exists := c.Get("CreateAbonementCommand")
	if !exists {
		logger.ErrorLogger.Printf("Cant find CreateAbonementCommand in context")
		return
	}

	createAbonementCommand, ok := createAbonementCommandAny.(*dtos.CreateAbonementCommand)
	if !ok {
		logger.ErrorLogger.Printf("CreateAbonementCommand has an invalid type")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "CreateAbonementCommand has an invalid type"})
		return
	}

	stream, err := (*h.abonementClient).CreateAbonement(context.Background())
	if err != nil {
		logger.ErrorLogger.Printf(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open CreateAbonement stream"})
		return
	}

	abonementDataForCreate := &abonementGRPC.AbonementDataForCreate{
		Title:        createAbonementCommand.Title,
		Validity:     createAbonementCommand.ValidityPeriod,
		VisitingTime: createAbonementCommand.VisitingTime,
		Price:        int32(createAbonementCommand.Price),
		ServicesIds:  createAbonementCommand.Services,
	}

	createAbonementRequestAbonementDataForCreate := &abonementGRPC.CreateAbonementRequest_AbonementDataForCreate{
		AbonementDataForCreate: abonementDataForCreate,
	}

	createAbonementRequest := &abonementGRPC.CreateAbonementRequest{
		Payload: createAbonementRequestAbonementDataForCreate,
	}

	err = stream.Send(createAbonementRequest)
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

			err = stream.Send(&abonementGRPC.CreateAbonementRequest{
				Payload: &abonementGRPC.CreateAbonementRequest_AbonementPhoto{
					AbonementPhoto: buffer[:n],
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
		"abonement": res.GetAbonementWithServices(),
	})
}

func (h *Handler) UpdateAbonement(c *gin.Context) {

	updateAbonementCommandAny, exists := c.Get("UpdateAbonementCommand")
	if !exists {
		logger.ErrorLogger.Printf("Cant find UpdateAbonementCommand in context")
		return
	}

	updateAbonementCommand, ok := updateAbonementCommandAny.(*dtos.UpdateAbonementCommand)
	if !ok {
		logger.ErrorLogger.Printf("UpdateAbonementCommand has an invalid type")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "UpdateAbonementCommand has an invalid type"})
		return
	}

	stream, err := (*h.abonementClient).UpdateAbonement(context.Background())
	if err != nil {
		fmt.Printf("failed to stat file: %v\n", err)
	}

	abonementDataForUpdate := &abonementGRPC.AbonementDataForUpdate{
		Id:           updateAbonementCommand.Id,
		Title:        updateAbonementCommand.Title,
		Validity:     updateAbonementCommand.ValidityPeriod,
		VisitingTime: updateAbonementCommand.VisitingTime,
		Price:        int32(updateAbonementCommand.Price),
		ServicesIds:  updateAbonementCommand.Services,
	}

	updateAbonementRequestAbonementDataForUpdate := &abonementGRPC.UpdateAbonementRequest_AbonementDataForUpdate{
		AbonementDataForUpdate: abonementDataForUpdate,
	}

	updateAbonementRequest := &abonementGRPC.UpdateAbonementRequest{
		Payload: updateAbonementRequestAbonementDataForUpdate,
	}

	err = stream.Send(updateAbonementRequest)
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

			err = stream.Send(&abonementGRPC.UpdateAbonementRequest{
				Payload: &abonementGRPC.UpdateAbonementRequest_AbonementPhoto{
					AbonementPhoto: buffer[:n],
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
		"abonement": res.GetAbonementWithServices(),
	})
}

func (h *Handler) DeleteAbonement(c *gin.Context) {

	id := c.Param("id")

	convertedId, err := uuid.Parse(id)
	if err != nil {

		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": fmt.Errorf("invalid id format")})
		return
	}

	_ = convertedId

	deleteAbonementByIdRequest := &abonementGRPC.DeleteAbonementByIdRequest{
		Id: id,
	}

	deletedAbonementRes, err := (*h.abonementClient).DeleteAbonementById(context.TODO(), deleteAbonementByIdRequest)
	if err != nil {
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"abonement": deletedAbonementRes.GetAbonementObject(),
	})
}
