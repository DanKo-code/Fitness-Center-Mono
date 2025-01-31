package usecase

import (
	"context"
	"github.com/DanKo-code/FitnessCenter-Coach/internal/dtos"
	"github.com/DanKo-code/FitnessCenter-Coach/internal/models"
	coachGRPC "github.com/DanKo-code/FitnessCenter-Protobuf/gen/FitnessCenter.protobuf.coach"
	"github.com/google/uuid"
)

type CoachUseCase interface {
	UpdateCoach(ctx context.Context, cmd *dtos.UpdateCoachCommand) (*models.Coach, error)
	CreateCoach(ctx context.Context, cmd *dtos.CreateCoachCommand) (*models.Coach, error)
	DeleteCoachById(ctx context.Context, id uuid.UUID) (*models.Coach, error)
	GetCoachById(ctx context.Context, uuid uuid.UUID) (*models.Coach, error)

	GetCoaches(ctx context.Context) ([]*models.Coach, error)
	GetCoachesWithServices(ctx context.Context) (*coachGRPC.GetCoachesWithServicesWithReviewsWithUsersResponse, error)
}
