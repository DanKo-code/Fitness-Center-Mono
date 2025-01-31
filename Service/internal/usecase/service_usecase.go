package usecase

import (
	"Service/internal/dtos"
	"Service/internal/models"
	"context"
	"github.com/google/uuid"
)

type ServiceUseCase interface {
	CreateService(ctx context.Context, cmd *dtos.CreateServiceCommand) (*models.Service, error)
	GetServiceById(ctx context.Context, id uuid.UUID) (*models.Service, error)
	UpdateService(ctx context.Context, cmd *dtos.UpdateServiceCommand) (*models.Service, error)
	DeleteServiceById(ctx context.Context, id uuid.UUID) (*models.Service, error)

	GetServices(ctx context.Context) ([]*models.Service, error)
	CreateCoachServices(ctx context.Context, cmd *dtos.CreateCoachServicesCommand) ([]*models.Service, error)
	CreateAbonemntServices(ctx context.Context, cmd *dtos.CreateAbonementServicesCommand) ([]*models.Service, error)
	GetAbonementsServices(ctx context.Context, ids []uuid.UUID) (map[uuid.UUID][]*models.Service, error)
	GetCoachesServices(ctx context.Context, ids []uuid.UUID) (map[uuid.UUID][]*models.Service, error)
	UpdateAbonementServices(ctx context.Context, abonementId uuid.UUID, servicesIds []uuid.UUID) ([]*models.Service, error)
	UpdateCoachServices(ctx context.Context, coachId uuid.UUID, servicesIds []uuid.UUID) ([]*models.Service, error)
}
