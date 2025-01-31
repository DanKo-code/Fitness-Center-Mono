package repository

import (
	"Service/internal/dtos"
	"Service/internal/models"
	"context"
	"github.com/google/uuid"
)

type ServiceRepository interface {
	CreateService(ctx context.Context, service *models.Service) error
	GetServiceById(ctx context.Context, id uuid.UUID) (*models.Service, error)
	UpdateService(ctx context.Context, cmd *dtos.UpdateServiceCommand) error
	DeleteService(ctx context.Context, id uuid.UUID) error

	GetServices(ctx context.Context) ([]*models.Service, error)
	CreateCoachServices(ctx context.Context, cmd *dtos.CreateCoachServicesCommand) error
	CreateAbonementServices(ctx context.Context, cmd *dtos.CreateAbonementServicesCommand) error
	GetServicesByIds(ctx context.Context, ids []uuid.UUID) ([]*models.Service, error)
	GetCoachServices(ctx context.Context, id uuid.UUID) ([]*models.Service, error)
	GetAbonementServices(ctx context.Context, id uuid.UUID) ([]*models.Service, error)
	GetAbonementsServices(ctx context.Context, ids []uuid.UUID) (map[uuid.UUID][]*models.Service, error)
	GetCoachesServices(ctx context.Context, ids []uuid.UUID) (map[uuid.UUID][]*models.Service, error)
	UpdateAbonementServices(ctx context.Context, abonementId uuid.UUID, servicesIds []uuid.UUID) error
	UpdateCoachServices(ctx context.Context, coachId uuid.UUID, servicesIds []uuid.UUID) error
}
