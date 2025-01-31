package grpc

import (
	"User/internal/dtos"
	customErrors "User/internal/errors"
	"User/internal/usecase"
	"User/pkg/logger"
	"context"
	"errors"
	userProtobuf "github.com/DanKo-code/FitnessCenter-Protobuf/gen/FitnessCenter.protobuf.user"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"io"
	"reflect"
	"strings"
	"time"
)

var _ userProtobuf.UserServer = (*UsergRPC)(nil)

type UsergRPC struct {
	userProtobuf.UnimplementedUserServer

	userUseCase  usecase.UserUseCase
	cloudUseCase usecase.CloudUseCase
}

func Register(
	gRPC *grpc.Server,
	userUseCase usecase.UserUseCase,
	cloudUseCase usecase.CloudUseCase,
) {
	userProtobuf.RegisterUserServer(gRPC, &UsergRPC{userUseCase: userUseCase, cloudUseCase: cloudUseCase})
}

func (u *UsergRPC) CreateUser(
	g grpc.ClientStreamingServer[
		userProtobuf.CreateUserRequest,
		userProtobuf.CreateUserResponse],
) error {

	userData, userPhoto, err := GetUserData(
		&g,
		func(chunk *userProtobuf.CreateUserRequest) interface{} {
			return chunk.GetUserDataForCreate()
		},
		func(chunk *userProtobuf.CreateUserRequest) []byte {
			return chunk.GetUserPhoto()
		},
	)
	if err != nil {
		return status.Error(codes.InvalidArgument, "invalid request data")
	}

	if userData == nil {
		logger.ErrorLogger.Printf("user data is empty")
		return status.Error(codes.InvalidArgument, "user data is empty")
	}

	castedUserData, ok := userData.(*userProtobuf.UserDataForCreate)
	if !ok {
		logger.ErrorLogger.Printf("user data is not of type UserProtobuf.UserDataForCreate")
		return status.Error(codes.InvalidArgument, "user data is not of type UserProtobuf.UserDataForCreate")
	}

	cmd := &dtos.CreateUserCommand{
		ID:       uuid.New(),
		Name:     castedUserData.Name,
		Email:    castedUserData.Email,
		Role:     castedUserData.Role,
		Password: castedUserData.Password,
	}

	var photoURL string
	randomID := uuid.New().String()
	if userPhoto != nil {
		url, err := u.cloudUseCase.PutObject(context.TODO(), userPhoto, "user/"+randomID)
		photoURL = url
		if err != nil {
			logger.ErrorLogger.Printf("Failed to create user photo in cloud: %v", err)
			return status.Error(codes.Internal, "Failed to create user photo in cloud")
		}
	}

	cmd.Photo = photoURL

	user, err := u.userUseCase.CreateUser(context.TODO(), cmd)
	if err != nil {
		return status.Error(codes.Internal, "Failed to create user")
	}

	userObject := &userProtobuf.UserObject{
		Id:          user.ID.String(),
		Email:       user.Email,
		Role:        user.Role,
		Photo:       user.Photo,
		Name:        user.Name,
		CreatedTime: user.CreatedTime.String(),
		UpdatedTime: user.UpdatedTime.String(),
	}

	response := &userProtobuf.CreateUserResponse{
		UserObject: userObject,
	}

	err = g.SendAndClose(response)
	if err != nil {
		logger.ErrorLogger.Printf("Failed to send user create response: %v", err)
		return status.Error(codes.Internal, "Failed to send user create response")
	}

	return nil
}

func (u *UsergRPC) GetUserById(ctx context.Context,
	request *userProtobuf.GetUserByIdRequest,
) (
	*userProtobuf.GetUserByIdResponse,
	error,
) {
	user, err := u.userUseCase.GetUserById(ctx, uuid.MustParse(request.Id))
	if err != nil {

		if errors.Is(err, customErrors.UserNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}

		return nil, err
	}

	userObject := &userProtobuf.UserObject{
		Id:          user.ID.String(),
		Email:       user.Email,
		Name:        user.Name,
		Role:        user.Role,
		Photo:       user.Photo,
		CreatedTime: user.CreatedTime.String(),
		UpdatedTime: user.UpdatedTime.String(),
	}

	response := &userProtobuf.GetUserByIdResponse{
		UserObject: userObject,
	}

	return response, nil
}

func (u *UsergRPC) UpdateUser(
	g grpc.ClientStreamingServer[
		userProtobuf.UpdateUserRequest,
		userProtobuf.UpdateUserResponse],
) error {

	userData, userPhoto, err := GetUserData(
		&g,
		func(chunk *userProtobuf.UpdateUserRequest) interface{} {
			return chunk.GetUserDataForUpdate()
		},
		func(chunk *userProtobuf.UpdateUserRequest) []byte {
			return chunk.GetUserPhoto()
		},
	)
	if err != nil {
		return status.Error(codes.InvalidArgument, "invalid request data")
	}

	if userData == nil {
		logger.ErrorLogger.Printf("user data is empty")
		return status.Error(codes.InvalidArgument, "user data is empty")
	}

	castedUserData, ok := userData.(*userProtobuf.UserDataForUpdate)
	if !ok {
		logger.ErrorLogger.Printf("user data is not of type UserProtobuf.UserDataForUpdate")
		return status.Error(codes.InvalidArgument, "user data is not of type UserProtobuf.UserDataForCreate")
	}

	cmd := &dtos.UpdateUserCommand{
		Id:          uuid.MustParse(castedUserData.Id),
		Name:        castedUserData.Name,
		Email:       castedUserData.Email,
		Role:        castedUserData.Role,
		UpdatedTime: time.Now(),
	}

	existingUser, err := u.userUseCase.GetUserById(context.TODO(), uuid.MustParse(castedUserData.Id))
	if err != nil {
		return status.Error(codes.NotFound, "user not found")
	}

	var photoURL string
	randomID := uuid.New().String()
	if userPhoto != nil {

		if existingUser.Photo != "" {
			prefix := "user/"
			index := strings.Index(existingUser.Photo, prefix)
			var s3PhotoKey string
			if index != -1 {
				s3PhotoKey = existingUser.Photo[index+len(prefix):]
			} else {
				logger.ErrorLogger.Printf("Prefix not found")
			}

			exists, err := u.cloudUseCase.ObjectExists(context.TODO(), "user/"+s3PhotoKey)
			if err != nil {
				return status.Error(codes.Internal, "can't find previous photo meta")
			}

			if exists {
				err := u.cloudUseCase.DeleteObject(context.TODO(), "user/"+s3PhotoKey)
				if err != nil {
					return err
				}
			}
		}

		url, err := u.cloudUseCase.PutObject(context.TODO(), userPhoto, "user/"+randomID)
		photoURL = url
		if err != nil {
			logger.ErrorLogger.Printf("Failed to create user photo in cloud: %v", err)
			return status.Error(codes.Internal, "Failed to create user photo in cloud")
		}
	}

	cmd.Photo = photoURL

	user, err := u.userUseCase.UpdateUser(context.TODO(), cmd)
	if err != nil {
		return status.Error(codes.Internal, "Failed to create user")
	}

	userObject := &userProtobuf.UserObject{
		Id:          user.ID.String(),
		Email:       user.Email,
		Role:        user.Role,
		Photo:       user.Photo,
		Name:        user.Name,
		CreatedTime: user.CreatedTime.String(),
		UpdatedTime: user.UpdatedTime.String(),
	}

	updateUserResponse := &userProtobuf.UpdateUserResponse{
		UserObject: userObject,
	}

	err = g.SendAndClose(updateUserResponse)
	if err != nil {
		logger.ErrorLogger.Printf("Failed to send user update response: %v", err)
		return err
	}

	return nil
}

func (u *UsergRPC) DeleteUserById(
	ctx context.Context,
	request *userProtobuf.DeleteUserByIdRequest,
) (*userProtobuf.DeleteUserByIdResponse,
	error,
) {
	deletedUser, err := u.userUseCase.DeleteUserById(ctx, uuid.MustParse(request.Id))
	if err != nil {
		return nil, err
	}

	logger.DebugLogger.Printf("deletedUser: %v\n", deletedUser)

	userObject := &userProtobuf.UserObject{
		Id:          deletedUser.ID.String(),
		Email:       deletedUser.Email,
		Role:        deletedUser.Role,
		Photo:       deletedUser.Photo,
		Name:        deletedUser.Name,
		CreatedTime: deletedUser.CreatedTime.String(),
		UpdatedTime: deletedUser.UpdatedTime.String(),
	}

	deleteUserByIdResponse := &userProtobuf.DeleteUserByIdResponse{
		UserObject: userObject,
	}

	return deleteUserByIdResponse, nil
}

func (u *UsergRPC) GetUserByEmail(
	ctx context.Context,
	request *userProtobuf.GetUserByEmailRequest,
) (*userProtobuf.GetUserByEmailResponse,
	error,
) {
	user, err := u.userUseCase.GetUserByEmail(ctx, request.Email)
	if err != nil {

		if errors.Is(err, customErrors.UserNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}

		return nil, err
	}

	userObject := &userProtobuf.UserObject{
		Id:          user.ID.String(),
		Email:       user.Email,
		Name:        user.Name,
		Role:        user.Role,
		Photo:       user.Photo,
		CreatedTime: user.CreatedTime.String(),
		UpdatedTime: user.UpdatedTime.String(),
	}

	response := &userProtobuf.GetUserByEmailResponse{
		UserObject: userObject,
	}

	return response, nil
}

func (u *UsergRPC) CheckPassword(
	ctx context.Context,
	request *userProtobuf.CheckPasswordRequest,
) (*emptypb.Empty,
	error,
) {

	checkPasswordCommand := &dtos.CheckPasswordCommand{
		UserId:   uuid.MustParse(request.UserId),
		Password: request.Password,
	}

	err := u.userUseCase.CheckPassword(ctx, checkPasswordCommand)
	if err != nil {

		if errors.Is(err, customErrors.UserNotFound) || errors.Is(err, customErrors.InvalidPassword) {
			return nil, status.Error(codes.Unauthenticated, err.Error())
		}

		return nil, status.Error(codes.Internal, err.Error())
	}

	return &emptypb.Empty{}, nil
}

func (u *UsergRPC) GetUsersByIds(
	ctx context.Context,
	request *userProtobuf.GetUsersByIdsRequest,
) (*userProtobuf.GetUsersByIdsResponse, error) {

	var ids []uuid.UUID
	for _, i2 := range request.UsersIds {
		ids = append(ids, uuid.MustParse(i2))
	}

	usersModels, err := u.userUseCase.GetUsersByIds(ctx, ids)
	if err != nil {
		return nil, err
	}

	response := &userProtobuf.GetUsersByIdsResponse{
		UsersObjects: nil,
	}

	for _, model := range usersModels {

		userObject := &userProtobuf.UserObject{
			Id:          model.ID.String(),
			Email:       model.Email,
			Role:        model.Role,
			Photo:       model.Photo,
			Name:        model.Name,
			CreatedTime: model.CreatedTime.String(),
			UpdatedTime: model.UpdatedTime.String(),
		}

		response.UsersObjects = append(response.UsersObjects, userObject)
	}

	return response, nil
}

func (u *UsergRPC) GetClients(ctx context.Context, empty *emptypb.Empty) (*userProtobuf.GetClientsResponse, error) {
	clients, err := u.userUseCase.GetClients(ctx)
	if err != nil {
		return nil, err
	}

	var clientObjects []*userProtobuf.UserObject
	for _, client := range clients {

		clientObject := &userProtobuf.UserObject{
			Id:          client.ID.String(),
			Email:       client.Email,
			Role:        client.Role,
			Photo:       client.Photo,
			Name:        client.Name,
			CreatedTime: client.CreatedTime.String(),
			UpdatedTime: client.UpdatedTime.String(),
		}

		clientObjects = append(clientObjects, clientObject)
	}

	response := &userProtobuf.GetClientsResponse{
		UserObjects: clientObjects,
	}

	return response, nil
}

func GetUserData[T any, R any](
	g *grpc.ClientStreamingServer[T, R],
	extractUserData func(chunk *T) interface{},
	extractUserPhoto func(chunk *T) []byte,
) (interface{},
	[]byte,
	error,
) {
	var userData interface{}
	var userPhoto []byte

	for {
		chunk, err := (*g).Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			logger.ErrorLogger.Printf("Error getting chunk: %v", err)
			return nil, nil, err
		}

		if ud := extractUserData(chunk); ud != nil && !reflect.ValueOf(ud).IsNil() {
			userData = ud
		}

		if uf := extractUserPhoto(chunk); uf != nil {
			userPhoto = append(userPhoto, uf...)
		}
	}

	return userData, userPhoto, nil
}
