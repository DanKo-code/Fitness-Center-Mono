package rest

import (
	"Gateway/internal/user/dtos"
	logger "Gateway/pkg/logger"
	"context"
	"fmt"
	userGRPC "github.com/DanKo-code/FitnessCenter-Protobuf/gen/FitnessCenter.protobuf.user"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"io"
	"net/http"
)

type Handler struct {
	userClient *grpc.ClientConn
	validator  *validator.Validate
}

func NewHandler(ssoClient *grpc.ClientConn, validator *validator.Validate) *Handler {
	return &Handler{
		userClient: ssoClient,
		validator:  validator,
	}
}

func (h *Handler) UpdateUser(c *gin.Context) {

	userId := c.Param("id")

	userIdFromToken, exists := c.Get("UserIdFromToken")
	if !exists {
		logger.ErrorLogger.Printf("Cant find UserIdFromToken in context")
		return
	}

	if userId != userIdFromToken {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Access denied: you cannot update another user's data",
		})
		return
	}

	cmd := &dtos.User{}

	form, err := c.MultipartForm()
	if err != nil {
		return
	}

	name, namOk := form.Value["name"]
	photo := form.File["photo"]

	if namOk {
		cmd.Name = name[0]
	}

	err = h.validator.Struct(cmd)
	if err != nil {
		logger.ErrorLogger.Printf("Error validating UpdateUserRequest: %v", err)

		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "имя должно содержать от 2 до 100 символов"})
		return
	}

	logger.InfoLogger.Printf("cmd: %v", cmd)

	userClient := userGRPC.NewUserClient(h.userClient)

	stream, err := userClient.UpdateUser(context.Background())
	if err != nil {
		fmt.Printf("failed to stat file: %v\n", err)
	}

	userDataForUpdate := &userGRPC.UserDataForUpdate{
		Id:    userId,
		Email: "",
		Name:  "",
		Role:  "",
	}

	userDataForUpdate.Id = userId
	if name != nil {
		userDataForUpdate.Name = name[0]
	}
	updateUserRequestUserDataForUpdate := &userGRPC.UpdateUserRequest_UserDataForUpdate{
		UserDataForUpdate: userDataForUpdate,
	}

	createUserRequest := &userGRPC.UpdateUserRequest{
		Payload: updateUserRequestUserDataForUpdate,
	}

	err = stream.Send(createUserRequest)
	if err != nil {
		return
	}

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

			err = stream.Send(&userGRPC.UpdateUserRequest{
				Payload: &userGRPC.UpdateUserRequest_UserPhoto{
					UserPhoto: buffer[:n],
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
		"user": res.GetUserObject(),
	})
}

func (h *Handler) GetClients(c *gin.Context) {

	userClient := userGRPC.NewUserClient(h.userClient)

	clients, err := userClient.GetClients(context.Background(), &emptypb.Empty{})
	if err != nil {
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"clients": clients.UserObjects,
	})
}

func (h *Handler) DeleteClientById(c *gin.Context) {

	userClient := userGRPC.NewUserClient(h.userClient)

	id := c.Param("id")

	convertedId, err := uuid.Parse(id)
	if err != nil {

		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": fmt.Errorf("invalid id format")})
		return
	}

	_ = convertedId

	deleteUserByIdRequest := &userGRPC.DeleteUserByIdRequest{
		Id: id,
	}

	client, err := userClient.DeleteUserById(context.Background(), deleteUserByIdRequest)
	if err != nil {
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"client": client.UserObject,
	})
}
