package usecase

import (
	"context"
	"github.com/DanKo-code/FitnessCenter-Review/internal/dtos"
	"github.com/DanKo-code/FitnessCenter-Review/internal/models"
	"github.com/google/uuid"
)

type ReviewUseCase interface {
	CreateCoachReview(ctx context.Context, cmd *dtos.CreateReviewCommand) (*models.Review, error)
	GetReviewById(ctx context.Context, id uuid.UUID) (*models.Review, error)
	UpdateReview(ctx context.Context, cmd *dtos.UpdateReviewCommand) (*models.Review, error)
	DeleteReviewById(ctx context.Context, id uuid.UUID) (*models.Review, error)

	GetCoachReviews(ctx context.Context, coachId uuid.UUID) ([]*models.Review, error)
	GetCoachesReviews(ctx context.Context, ids []uuid.UUID) (map[uuid.UUID][]*models.Review, error)
}
