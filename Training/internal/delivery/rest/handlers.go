package rest

import (
	"Training/internal/model"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"time"
)

const (
	BookTrainingStatus = "booked"
)

type trainingUseCase interface {
	Insert(ctx context.Context, training model.Training) (model.Training, error)
}

type Handlers struct {
	useCase trainingUseCase
}

func NewHandler(useCase trainingUseCase) Handlers {
	return Handlers{
		useCase,
	}
}

func (h Handlers) Insert(c *gin.Context) {

	var bookTrainingQuery BookTrainingQuery

	err := c.ShouldBindJSON(&bookTrainingQuery)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	trainingModel := model.Training{
		Id:          uuid.New(),
		TimeFrom:    bookTrainingQuery.TimeFrom,
		TimeUntil:   bookTrainingQuery.TimeUntil,
		Status:      BookTrainingStatus,
		CoachId:     bookTrainingQuery.CoachId,
		ClientId:    bookTrainingQuery.ClientId,
		CreatedTime: time.Now(),
		UpdatedTime: time.Now(),
	}

	insertedTrainingModel, err := h.useCase.Insert(context.Background(), trainingModel)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, insertedTrainingModel)
}
