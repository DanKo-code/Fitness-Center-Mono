package user_usecase

import (
	"context"
	"fmt"
	"github.com/DanKo-code/FitnessCenter-Coach/internal/dtos"
	customErrors "github.com/DanKo-code/FitnessCenter-Coach/internal/errors"
	"github.com/DanKo-code/FitnessCenter-Coach/internal/models"
	"github.com/DanKo-code/FitnessCenter-Coach/internal/repository"
	"github.com/DanKo-code/FitnessCenter-Coach/internal/usecase"
	"github.com/DanKo-code/FitnessCenter-Coach/pkg/logger"
	coachGRPC "github.com/DanKo-code/FitnessCenter-Protobuf/gen/FitnessCenter.protobuf.coach"
	reviewGRPC "github.com/DanKo-code/FitnessCenter-Protobuf/gen/FitnessCenter.protobuf.review"
	serviceGRPC "github.com/DanKo-code/FitnessCenter-Protobuf/gen/FitnessCenter.protobuf.service"
	userGRPC "github.com/DanKo-code/FitnessCenter-Protobuf/gen/FitnessCenter.protobuf.user"
	"github.com/google/uuid"
	"strings"
	"time"
)

type CoachUseCase struct {
	coachRepo     repository.CoachRepository
	serviceClient *serviceGRPC.ServiceClient
	reviewClient  *reviewGRPC.ReviewClient
	userClient    *userGRPC.UserClient
	cloudUseCase  usecase.CloudUseCase
}

func NewCoachUseCase(
	coachRepo repository.CoachRepository,
	serviceClient *serviceGRPC.ServiceClient,
	reviewClient *reviewGRPC.ReviewClient,
	userClient *userGRPC.UserClient,
	cloudUseCase usecase.CloudUseCase,
) *CoachUseCase {
	return &CoachUseCase{
		coachRepo:     coachRepo,
		serviceClient: serviceClient,
		reviewClient:  reviewClient,
		userClient:    userClient,
		cloudUseCase:  cloudUseCase,
	}
}

func (c *CoachUseCase) CreateCoach(
	ctx context.Context,
	cmd *dtos.CreateCoachCommand,
) (*models.Coach, error) {

	coach := &models.Coach{
		Id:          cmd.Id,
		Name:        cmd.Name,
		Description: cmd.Description,
		Photo:       cmd.Photo,
		UpdatedTime: time.Now(),
		CreatedTime: time.Now(),
	}

	createdCoach, err := c.coachRepo.CreateCoach(ctx, coach)
	if err != nil {
		return nil, err
	}

	return createdCoach, nil
}

func (c *CoachUseCase) GetCoachById(
	ctx context.Context,
	id uuid.UUID,
) (*models.Coach, error) {
	coach, err := c.coachRepo.GetCoachById(ctx, id)
	if err != nil {
		return nil, err
	}

	return coach, nil
}

func (c *CoachUseCase) UpdateCoach(
	ctx context.Context,
	cmd *dtos.UpdateCoachCommand,
) (*models.Coach, error) {

	err := c.coachRepo.UpdateCoach(ctx, cmd)
	if err != nil {
		return nil, err
	}

	coach, err := c.coachRepo.GetCoachById(ctx, cmd.Id)
	if err != nil {
		return nil, err
	}

	return coach, nil
}

func (c *CoachUseCase) DeleteCoachById(
	ctx context.Context,
	id uuid.UUID,
) (*models.Coach, error) {
	coach, err := c.coachRepo.GetCoachById(ctx, id)
	if err != nil {
		return nil, customErrors.CoachNotFound
	}

	err = c.coachRepo.DeleteCoachById(ctx, id)
	if err != nil {
		return nil, err
	}

	if coach.Photo != "" {
		prefix := "coach/"
		index := strings.Index(coach.Photo, prefix)
		var s3PhotoKey string
		if index != -1 {
			s3PhotoKey = coach.Photo[index+len(prefix):]
		} else {
			logger.ErrorLogger.Printf("Prefix not found")
			return nil, fmt.Errorf("prefix not found")
		}
		err = c.cloudUseCase.DeleteObject(ctx, "coach/"+s3PhotoKey)
		if err != nil {
			return nil, err
		}
	}
	
	return coach, nil
}

func (c *CoachUseCase) GetCoaches(
	ctx context.Context,
) ([]*models.Coach, error) {

	coaches, err := c.coachRepo.GetCoaches(ctx)
	if err != nil {
		return nil, err
	}

	return coaches, nil
}

func (c *CoachUseCase) GetCoachesWithServices(
	ctx context.Context,
) (*coachGRPC.GetCoachesWithServicesWithReviewsWithUsersResponse, error) {
	coaches, err := c.coachRepo.GetCoaches(ctx)
	if err != nil {
		logger.ErrorLogger.Printf("Failed GetCoaches: %s", err)
		return nil, err
	}

	//get coaches ids
	var coachesIds []string
	for _, coach := range coaches {
		coachesIds = append(coachesIds, coach.Id.String())
	}

	//get CoachesIdsWithServices
	getCoachesServicesRequest := &serviceGRPC.GetCoachesServicesRequest{
		CoachIds: coachesIds,
	}
	getCoachesServicesResponse, err := (*c.serviceClient).GetCoachesServices(ctx, getCoachesServicesRequest)
	if err != nil {
		logger.ErrorLogger.Printf("Failed GetCoachesServices: %s", err)
		return nil, err
	}
	coachIdServices := make(map[string][]*serviceGRPC.ServiceObject)
	for _, i2 := range getCoachesServicesResponse.CoachIdsWithServices {
		coachIdServices[i2.CoachId] = append(coachIdServices[i2.CoachId], i2.ServiceObjects...)
	}

	//get CoachesIdsWithReviews
	getCoachesReviewsRequest := &reviewGRPC.GetCoachesReviewsRequest{
		CoachesIds: coachesIds,
	}
	GetCoachesReviewsResponse, err := (*c.reviewClient).GetCoachesReviews(ctx, getCoachesReviewsRequest)
	if err != nil {
		return nil, err
	}
	coachIdReviews := make(map[string][]*reviewGRPC.ReviewObject)
	for _, i2 := range GetCoachesReviewsResponse.CoachIdWithReviewObject {
		coachIdReviews[i2.CoachId] = append(coachIdReviews[i2.CoachId], i2.ReviewObjects...)
	}

	//get unique user ids
	uniqUserIds := make(map[string]struct{})
	for _, object := range GetCoachesReviewsResponse.CoachIdWithReviewObject {

		for _, reviewObject := range object.ReviewObjects {
			if _, ok := uniqUserIds[reviewObject.UserId]; !ok {
				uniqUserIds[reviewObject.UserId] = struct{}{}
			}
		}
	}
	var uniqUserIdsSl []string
	for key, _ := range uniqUserIds {
		uniqUserIdsSl = append(uniqUserIdsSl, key)
	}

	//get userIdsWithUsers
	getUsersByIdsRequest := &userGRPC.GetUsersByIdsRequest{
		UsersIds: uniqUserIdsSl,
	}
	getUsersByIdsResponse, err := (*c.userClient).GetUsersByIds(ctx, getUsersByIdsRequest)
	if err != nil {
		return nil, err
	}
	userIdUser := make(map[string]*userGRPC.UserObject)
	for _, i2 := range getUsersByIdsResponse.UsersObjects {
		userIdUser[i2.Id] = i2
	}

	// create method response
	var coachWithServicesWithReviewsWithUsersSl []*coachGRPC.CoachWithServicesWithReviewsWithUsers

	response := &coachGRPC.GetCoachesWithServicesWithReviewsWithUsersResponse{
		CoachWithServicesWithReviewsWithUsers: nil,
	}

	for _, coach := range coaches {

		coachObject := &coachGRPC.CoachObject{
			Id:          coach.Id.String(),
			Name:        coach.Name,
			Description: coach.Description,
			Photo:       coach.Photo,
			CreatedTime: coach.CreatedTime.String(),
			UpdatedTime: coach.UpdatedTime.String(),
		}

		coachServices := coachIdServices[coach.Id.String()]

		var reviewsWithUsers []*coachGRPC.ReviewWithUser
		for _, i2 := range coachIdReviews[coach.Id.String()] {
			reviewWithUser := &coachGRPC.ReviewWithUser{
				ReviewObject: i2,
				UserObject:   userIdUser[i2.UserId],
			}

			reviewsWithUsers = append(reviewsWithUsers, reviewWithUser)
		}

		coachWithServicesWithReviewsWithUsers := &coachGRPC.CoachWithServicesWithReviewsWithUsers{
			Coach:          coachObject,
			Services:       coachServices,
			ReviewWithUser: reviewsWithUsers,
		}

		coachWithServicesWithReviewsWithUsersSl = append(coachWithServicesWithReviewsWithUsersSl, coachWithServicesWithReviewsWithUsers)
	}

	response.CoachWithServicesWithReviewsWithUsers = coachWithServicesWithReviewsWithUsersSl

	return response, nil
}
