package usecase

import (
	"User/internal/dtos"
	"User/internal/models"
	"context"
	"github.com/google/uuid"
)

type UserUseCase interface {
	CreateUser(ctx context.Context, cmd *dtos.CreateUserCommand) (*models.User, error)
	GetUserById(ctx context.Context, id uuid.UUID) (*models.User, error)
	UpdateUser(ctx context.Context, cmd *dtos.UpdateUserCommand) (*models.User, error)
	DeleteUserById(ctx context.Context, id uuid.UUID) (*models.User, error)

	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	CheckPassword(ctx context.Context, cmd *dtos.CheckPasswordCommand) error
	GetUsersByIds(ctx context.Context, ids []uuid.UUID) ([]*models.User, error)
	GetClients(ctx context.Context) ([]*models.User, error)
	GetAdmins(ctx context.Context) ([]*models.User, error)
}
