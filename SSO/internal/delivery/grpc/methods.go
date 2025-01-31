package grpc

import (
	"SSO/internal/dtos"
	"SSO/internal/ssoErrors"
	"SSO/internal/usecase"
	"context"
	"errors"
	ssoGRPC "github.com/DanKo-code/FitnessCenter-Protobuf/gen/FitnessCenter.protobuf.sso"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type SSOgRPC struct {
	ssoGRPC.UnimplementedSSOServer
	useCase usecase.UseCase
}

func Register(gRPC *grpc.Server, useCase usecase.UseCase) {
	ssoGRPC.RegisterSSOServer(gRPC, &SSOgRPC{useCase: useCase})
}

func (s SSOgRPC) SignUp(ctx context.Context, request *ssoGRPC.SignUpRequest) (*ssoGRPC.SignUpResponse, error) {

	signUpRequestCommand := &dtos.SignUpRequestCommand{
		Name:        request.Name,
		Email:       request.Email,
		Password:    request.Password,
		FingerPrint: request.FingerPrint,
	}

	signUpResponseCommand, err := s.useCase.SignUp(ctx, signUpRequestCommand)
	if err != nil {
		if errors.Is(err, ssoErrors.UserAlreadyExists) {
			return nil, status.Error(codes.AlreadyExists, err.Error())
		}

		return nil, status.Error(codes.Internal, err.Error())
	}

	signUpResponse := &ssoGRPC.SignUpResponse{
		AccessToken:            signUpResponseCommand.AccessToken,
		RefreshToken:           signUpResponseCommand.RefreshToken,
		AccessTokenExpiration:  signUpResponseCommand.AccessTokenExpiration,
		RefreshTokenExpiration: signUpResponseCommand.RefreshTokenExpiration,
		User:                   signUpResponseCommand.User,
	}

	return signUpResponse, nil
}

func (s SSOgRPC) SignIn(ctx context.Context, request *ssoGRPC.SignInRequest) (*ssoGRPC.SignInResponse, error) {

	signInRequestCommand := &dtos.SignInRequestCommand{
		Email:       request.Email,
		Password:    request.Password,
		FingerPrint: request.FingerPrint,
	}

	signInResponseCommand, err := s.useCase.SignIn(ctx, signInRequestCommand)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	response := &ssoGRPC.SignInResponse{
		AccessToken:            signInResponseCommand.AccessToken,
		RefreshToken:           signInResponseCommand.RefreshToken,
		AccessTokenExpiration:  signInResponseCommand.AccessTokenExpiration,
		RefreshTokenExpiration: signInResponseCommand.RefreshTokenExpiration,
		User:                   signInResponseCommand.User,
	}

	return response, nil
}

func (s SSOgRPC) LogOut(ctx context.Context, request *ssoGRPC.LogOutRequest) (*emptypb.Empty, error) {

	refreshToken := request.GetRefreshToken()

	err := s.useCase.LogOut(ctx, refreshToken)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &emptypb.Empty{}, nil
}

func (s SSOgRPC) Refresh(ctx context.Context, request *ssoGRPC.RefreshRequest) (*ssoGRPC.RefreshResponse, error) {

	refreshRequestCommand := &dtos.RefreshRequestCommand{
		FingerPrint:  request.GetFingerPrint(),
		RefreshToken: request.GetRefreshToken(),
	}

	refreshResponseCommand, err := s.useCase.Refresh(ctx, refreshRequestCommand)
	if err != nil {
		if errors.Is(err, ssoErrors.InvalidFingerPrint) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	response := &ssoGRPC.RefreshResponse{
		AccessToken:            refreshResponseCommand.AccessToken,
		RefreshToken:           refreshResponseCommand.RefreshToken,
		AccessTokenExpiration:  refreshResponseCommand.AccessTokenExpiration,
		RefreshTokenExpiration: refreshResponseCommand.RefreshTokenExpiration,
		User:                   refreshResponseCommand.User,
	}

	return response, nil
}
