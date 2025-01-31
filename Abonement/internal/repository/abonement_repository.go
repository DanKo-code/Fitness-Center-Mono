package repository

import (
	"context"
	"github.com/DanKo-code/Fitness-Center-Abonement/internal/dtos"
	"github.com/DanKo-code/Fitness-Center-Abonement/internal/models"
	"github.com/google/uuid"
)

type AbonementRepository interface {
	CreateAbonement(ctx context.Context, abonement *models.Abonement) (*models.Abonement, error)
	GetAbonementById(ctx context.Context, id uuid.UUID) (*models.Abonement, error)
	UpdateAbonement(ctx context.Context, cmd *dtos.UpdateAbonementCommand) error
	DeleteAbonementById(ctx context.Context, id uuid.UUID) error

	GetAbonementes(ctx context.Context) ([]*models.Abonement, error)
	GetAbonementsByIds(ctx context.Context, ids []uuid.UUID) ([]*models.Abonement, error)
}
