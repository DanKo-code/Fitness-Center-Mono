package review_usecase

import (
	"context"
	coachGRPC "github.com/DanKo-code/FitnessCenter-Protobuf/gen/FitnessCenter.protobuf.coach"
	userGRPC "github.com/DanKo-code/FitnessCenter-Protobuf/gen/FitnessCenter.protobuf.user"
	"github.com/DanKo-code/FitnessCenter-Review/internal/dtos"
	customErrors "github.com/DanKo-code/FitnessCenter-Review/internal/errors"
	"github.com/DanKo-code/FitnessCenter-Review/internal/models"
	"github.com/DanKo-code/FitnessCenter-Review/internal/repository"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

type ReviewUseCase struct {
	reviewRepo  repository.ReviewRepository
	coachClient *coachGRPC.CoachClient
	userClient  *userGRPC.UserClient
}

func NewReviewUseCase(
	reviewRepo repository.ReviewRepository,
	coachClient *coachGRPC.CoachClient,
	userClient *userGRPC.UserClient,
) *ReviewUseCase {
	return &ReviewUseCase{
		reviewRepo:  reviewRepo,
		coachClient: coachClient,
		userClient:  userClient,
	}
}

func (u *ReviewUseCase) CreateCoachReview(ctx context.Context, cmd *dtos.CreateReviewCommand) (*models.Review, error) {

	getUserByIdRequest := &userGRPC.GetUserByIdRequest{Id: cmd.UserId.String()}
	_, err := (*u.userClient).GetUserById(ctx, getUserByIdRequest)
	if err != nil {

		st, ok := status.FromError(err)

		if !ok {
			return nil, nil
		}

		switch st.Code() {
		case codes.NotFound:
			return nil, customErrors.UserNotFound
		default:
			return nil, customErrors.InternalUserServerError
		}
	}

	getCoachByIdRequest := &coachGRPC.GetCoachByIdRequest{Id: cmd.CoachId.String()}
	_, err = (*u.coachClient).GetCoachById(ctx, getCoachByIdRequest)
	if err != nil {

		st, ok := status.FromError(err)

		if !ok {
			return nil, nil
		}

		switch st.Code() {
		case codes.NotFound:
			return nil, customErrors.CoachNotFound
		default:
			return nil, customErrors.InternalCoachServerError
		}
	}

	review := &models.Review{
		Id:          uuid.New(),
		UserId:      cmd.UserId,
		Body:        cmd.Body,
		UpdatedTime: time.Now(),
		CreatedTime: time.Now(),
	}

	err = u.reviewRepo.CreateCoachReview(ctx, review, cmd.CoachId)
	if err != nil {
		return nil, err
	}

	return review, nil
}

func (u *ReviewUseCase) GetReviewById(ctx context.Context, id uuid.UUID) (*models.Review, error) {
	review, err := u.reviewRepo.GetReviewById(ctx, id)
	if err != nil {
		return nil, err
	}

	return review, nil
}

func (u *ReviewUseCase) UpdateReview(ctx context.Context, cmd *dtos.UpdateReviewCommand) (*models.Review, error) {

	cmd.UpdatedTime = time.Now()

	err := u.reviewRepo.UpdateReview(ctx, cmd)
	if err != nil {
		return nil, err
	}

	review, err := u.reviewRepo.GetReviewById(ctx, cmd.Id)
	if err != nil {
		return nil, err
	}

	return review, nil
}

func (u *ReviewUseCase) DeleteReviewById(ctx context.Context, id uuid.UUID) (*models.Review, error) {

	review, err := u.reviewRepo.GetReviewById(ctx, id)
	if err != nil {
		return nil, err
	}

	err = u.reviewRepo.DeleteReviewById(ctx, id)
	if err != nil {
		return nil, err
	}

	return review, nil
}

func (u *ReviewUseCase) GetCoachReviews(ctx context.Context, coachId uuid.UUID) ([]*models.Review, error) {
	reviews, err := u.reviewRepo.GetCoachReviews(ctx, coachId)
	if err != nil {
		return nil, err
	}

	return reviews, nil
}

func (u *ReviewUseCase) GetCoachesReviews(ctx context.Context, ids []uuid.UUID) (map[uuid.UUID][]*models.Review, error) {
	reviews, err := u.reviewRepo.GetCoachesReviews(ctx, ids)
	if err != nil {
		return nil, err
	}

	return reviews, nil
}
