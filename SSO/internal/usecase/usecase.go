package usecase

import (
	"SSO/internal/dtos"
	"context"
)

type UseCase interface {
	SignUp(ctx context.Context, cmd *dtos.SignUpRequestCommand) (*dtos.SignUpResponseCommand, error)
	SignIn(ctx context.Context, siReq *dtos.SignInRequestCommand) (*dtos.SignInResponseCommand, error)
	LogOut(ctx context.Context, refreshToken string) error
	Refresh(ctx context.Context, cmd *dtos.RefreshRequestCommand) (*dtos.RefreshResponseCommand, error)
}
