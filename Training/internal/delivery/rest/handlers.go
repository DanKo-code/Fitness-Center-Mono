package rest

import (
	"Training/internal/model"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	BookTrainingStatus = "booked"
)

type trainingUseCase interface {
	Insert(ctx context.Context, training model.Training) (model.Training, error)
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
	roomId := c.Param("roomId")
	/*userIdFromToken, exists := c.Get("UserIdFromToken")
	if !exists {
		c.Error(fmt.Errorf("cant find UserIdFromToken in context"))
		return
	}

	var joinTrainingQuery JoinTrainingQuery
	err := c.ShouldBindJSON(&joinTrainingQuery)
	if err != nil {
		c.Error(fmt.Errorf("cant Bind JoinTrainingQuery"))
		return
	}*/

	coachId := c.Query("coachId") // ?userId=123
	token := c.Query("token")     // ?token=abcd

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
	roomKey.ClientId = uuid.MustParse(claims.UserId)
	roomKey.CoachId = uuid.MustParse(coachId)

	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.Error(err)
		return
	}

	h.roomMap.InsertIntoRoom(roomKey, ws)
	for {
		participants := h.roomMap.Get(roomKey)
		if participants == nil || len(participants) == 0 {
			return
		}

		var msg model.BroadcastMsg

		log.Println("Start ws.ReadJSON: ", msg)
		err := ws.ReadJSON(&msg.Message)
		if err != nil {
			log.Fatal("Read Error: ", err)
		}

		msg.Client = ws
		msg.RoomKey = roomKey

		log.Println("End ws.ReadJSON: ", msg)

		h.broadcast <- msg
	}
}
