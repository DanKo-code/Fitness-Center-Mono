package repository

import (
	"context"
	"github.com/DanKo-code/FitnessCenter-Coach/internal/dtos"
	"github.com/DanKo-code/FitnessCenter-Coach/internal/models"
	"github.com/google/uuid"
)

type CoachRepository interface {
	CreateCoach(ctx context.Context, coach *models.Coach) (*models.Coach, error)
	GetCoachById(ctx context.Context, id uuid.UUID) (*models.Coach, error)
	UpdateCoach(ctx context.Context, cmd *dtos.UpdateCoachCommand) error
	DeleteCoachById(ctx context.Context, id uuid.UUID) error

	GetCoaches(ctx context.Context) ([]*models.Coach, error)
}
