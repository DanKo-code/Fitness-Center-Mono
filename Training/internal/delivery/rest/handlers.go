package rest

import (
	"Training/internal/model"
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

const (
	BookTrainingStatus = "booked"
)

type trainingUseCase interface {
	Insert(ctx context.Context, training model.Training) (model.Training, error)
	GetTrainingsByDateAndCoach(ctx context.Context, date string, coachId string) ([]model.Training, error)
	IsValidParticipant(ctx context.Context, roomId, clientId, coachId uuid.UUID, role string) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetActiveTrainingsByClient(ctx context.Context, clientId string) ([]model.TrainingWithCoachDetails, error)
}

type Handlers struct {
	useCase   trainingUseCase
	roomMap   *model.RoomMap
	broadcast chan model.BroadcastMsg
}

func NewHandler(useCase trainingUseCase, roomMap *model.RoomMap, broadcast chan model.BroadcastMsg) Handlers {
	return Handlers{
		useCase,
		roomMap,
		broadcast,
	}
}

func (h Handlers) GetActiveTrainingsByCoachId(c *gin.Context) {
	clientId := c.Param("clientId")

	slog.Info("clientId", clientId)

	trainingsWithCoaches, err := h.useCase.GetActiveTrainingsByClient(context.Background(), clientId)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"trainingsWithCoaches": trainingsWithCoaches})
}

func (h Handlers) Insert(c *gin.Context) {

	var bookTrainingQuery BookTrainingQuery

	err := c.ShouldBindJSON(&bookTrainingQuery)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	convertedTimeFrom, convertedTimeUntil, err := createTimeFromAndTimeUntil(bookTrainingQuery.Date, bookTrainingQuery.TimeFrom, bookTrainingQuery.TimeUntil)

	userId, ok := c.Get("UserIdFromToken")
	if !ok {
		c.Error(fmt.Errorf("UserIdFromToken not dound in context"))
		return
	}

	trainingModel := model.Training{
		Id:          uuid.New(),
		TimeFrom:    convertedTimeFrom,
		TimeUntil:   convertedTimeUntil,
		Status:      BookTrainingStatus,
		CoachId:     bookTrainingQuery.CoachId,
		ClientId:    uuid.MustParse(userId.(string)),
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

func (h Handlers) Delete(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		slog.Info("id not setted")
		c.JSON(http.StatusBadRequest, gin.H{"error": "id not setted"})
		return
	}

	idUUID, err := uuid.Parse(id)
	if err != nil {
		slog.Info("id not uuid")
		c.JSON(http.StatusBadRequest, gin.H{"error": "id not uuid"})
		return
	}

	err = h.useCase.Delete(c.Request.Context(), idUUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusOK)
}

func createTimeFromAndTimeUntil(date, timeFrom, timeUntil string) (convertedTimeFrom, convertedTimeUntil time.Time, err error) {
	parsedDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	parsedTimeFrom, err := time.Parse("15:04", timeFrom)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	parsedTimeUntil, err := time.Parse("15:04", timeUntil)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	convertedTimeFrom = time.Date(
		parsedDate.Year(), parsedDate.Month(), parsedDate.Day(),
		parsedTimeFrom.Hour(), parsedTimeFrom.Minute(), 0, 0, time.Local,
	)

	convertedTimeUntil = time.Date(
		parsedDate.Year(), parsedDate.Month(), parsedDate.Day(),
		parsedTimeUntil.Hour(), parsedTimeUntil.Minute(), 0, 0, time.Local,
	)

	return convertedTimeFrom, convertedTimeUntil, nil
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Claims struct {
	UserId string `json:"user_id"`
	Role   string `json:"role"`
	Exp    int64  `json:"exp"`
	jwt.RegisteredClaims
}

func VerifyAccessToken(accessToken string) (*Claims, error) {
	secretKey := os.Getenv("JWT_SECRET")
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(accessToken, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}

func (h Handlers) Join(c *gin.Context) {

	slog.Info("Join")

	roomId := c.Param("roomId")

	coachId := c.Query("coachId")
	token := c.Query("token")

	slog.Info("before VerifyAccessToken")

	claims, err := VerifyAccessToken(token)
	if err != nil {
		return
	}

	roomKey := model.RoomMapKey{
		RoomId:   uuid.UUID{},
		ClientId: uuid.UUID{},
		CoachId:  uuid.UUID{},
	}

	roomKey.RoomId = uuid.MustParse(roomId)
	roomKey.ClientId = uuid.Nil
	roomKey.CoachId = uuid.Nil

	slog.Info("before IsValidParticipant")

	err = h.useCase.IsValidParticipant(c.Request.Context(), roomKey.RoomId, uuid.MustParse(claims.UserId), uuid.MustParse(coachId), claims.Role)
	if err != nil {
		slog.Info("error IsValidParticipant", err)
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	slog.Info("before upgrader.Upgrade")

	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.Error(err)
		return
	}

	slog.Info("after upgrader.Upgrade")

	h.roomMap.InsertIntoRoom(roomKey, ws)
	for {
		participants := h.roomMap.Get(roomKey)
		if participants == nil || len(participants) == 0 {
			return
		}

		var msg model.BroadcastMsg

		err := ws.ReadJSON(&msg.Message)
		if err != nil {
			log.Printf("Read Error: ", err)
			h.roomMap.DeleteFromRoom(roomKey, ws)
			return
		}

		msg.Client = ws
		msg.RoomKey = roomKey

		h.broadcast <- msg
	}
}

func (h Handlers) GetTrainingsByDayAndCoach(c *gin.Context) {

	slog.Info("GetTrainingsByDayAndCoach")

	day := c.Param("day")
	coachId := c.Param("coachId")

	_, err := time.Parse("2006-01-02", day)
	if err != nil {
		c.Error(err)
		return
	}

	trainings, err := h.useCase.GetTrainingsByDateAndCoach(context.Background(), day, coachId)
	if err != nil {
		c.Error(err)
		return
	}

	if trainings == nil {
		trainings = []model.Training{}
	}

	c.JSON(http.StatusOK, gin.H{"trainings": trainings})
}
