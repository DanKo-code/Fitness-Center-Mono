package usecase

import (
	"context"
	"github.com/DanKo-code/Fitness-Center-Abonement/internal/dtos"
	"github.com/DanKo-code/Fitness-Center-Abonement/internal/models"
	"github.com/google/uuid"
)

type AbonementUseCase interface {
	UpdateAbonement(ctx context.Context, cmd *dtos.UpdateAbonementCommand) (*models.Abonement, error)
	CreateAbonement(ctx context.Context, cmd *dtos.CreateAbonementCommand) (*models.Abonement, error)
	DeleteAbonementById(ctx context.Context, id uuid.UUID) (*models.Abonement, error)
	GetAbonementById(ctx context.Context, uuid uuid.UUID) (*models.Abonement, error)

	GetAbonementes(ctx context.Context) ([]*models.Abonement, error)
	GetAbonementsWithServices(ctx context.Context) ([]*dtos.AbonementWithServices, error)
	GetAbonementsByIds(ctx context.Context, ids []uuid.UUID) ([]*models.Abonement, error)
}
