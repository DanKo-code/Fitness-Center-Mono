package grpc

import (
	"context"
	"errors"
	reviewProtobuf "github.com/DanKo-code/FitnessCenter-Protobuf/gen/FitnessCenter.protobuf.review"
	"github.com/DanKo-code/FitnessCenter-Review/internal/dtos"
	customErrors "github.com/DanKo-code/FitnessCenter-Review/internal/errors"
	"github.com/DanKo-code/FitnessCenter-Review/internal/usecase"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ReviewgRPC struct {
	reviewProtobuf.UnimplementedReviewServer

	reviewUseCase usecase.ReviewUseCase
}

func Register(gRPC *grpc.Server, reviewUseCase usecase.ReviewUseCase) {
	reviewProtobuf.RegisterReviewServer(gRPC, &ReviewgRPC{reviewUseCase: reviewUseCase})
}

func (r *ReviewgRPC) CreateCoachReview(
	ctx context.Context,
	request *reviewProtobuf.CreateCoachReviewRequest,
) (*reviewProtobuf.CreateCoachReviewResponse, error) {

	createReviewCommand := &dtos.CreateReviewCommand{
		UserId:  uuid.MustParse(request.ReviewDataForCreate.UserId),
		Body:    request.ReviewDataForCreate.Body,
		CoachId: uuid.MustParse(request.ReviewDataForCreate.CoachId),
	}

	review, err := r.reviewUseCase.CreateCoachReview(ctx, createReviewCommand)
	if err != nil {
		return nil, err
	}

	reviewObject := &reviewProtobuf.ReviewObject{
		Id:          review.Id.String(),
		Body:        review.Body,
		UserId:      review.UserId.String(),
		CreatedTime: review.CreatedTime.String(),
		UpdatedTime: review.UpdatedTime.String(),
	}

	response := &reviewProtobuf.CreateCoachReviewResponse{
		ReviewObject: reviewObject,
	}

	return response, nil
}

func (r *ReviewgRPC) GetReviewById(
	ctx context.Context,
	request *reviewProtobuf.GetReviewByIdRequest,
) (*reviewProtobuf.GetReviewByIdResponse, error) {

	review, err := r.reviewUseCase.GetReviewById(ctx, uuid.MustParse(request.Id))
	if err != nil {

		if errors.Is(err, customErrors.ReviewNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}

		return nil, err
	}

	reviewObject := &reviewProtobuf.ReviewObject{
		Id:          review.Id.String(),
		UserId:      review.UserId.String(),
		Body:        review.Body,
		CreatedTime: review.CreatedTime.String(),
		UpdatedTime: review.UpdatedTime.String(),
	}

	response := &reviewProtobuf.GetReviewByIdResponse{
		ReviewObject: reviewObject,
	}

	return response, nil
}

func (r *ReviewgRPC) UpdateReview(
	ctx context.Context,
	request *reviewProtobuf.UpdateReviewRequest,
) (*reviewProtobuf.UpdateReviewResponse, error) {

	updateReviewCommand := &dtos.UpdateReviewCommand{
		Id:   uuid.MustParse(request.ReviewDataForUpdate.Id),
		Body: request.ReviewDataForUpdate.Body,
	}

	updatedReview, err := r.reviewUseCase.UpdateReview(ctx, updateReviewCommand)
	if err != nil {
		return nil, err
	}

	reviewObject := &reviewProtobuf.ReviewObject{
		Id:          updatedReview.Id.String(),
		Body:        updatedReview.Body,
		CreatedTime: updatedReview.CreatedTime.String(),
		UpdatedTime: updatedReview.UpdatedTime.String(),
	}

	response := &reviewProtobuf.UpdateReviewResponse{
		ReviewObject: reviewObject,
	}

	return response, nil
}

func (r *ReviewgRPC) DeleteReviewById(
	ctx context.Context,
	request *reviewProtobuf.DeleteReviewByIdRequest,
) (*reviewProtobuf.DeleteReviewByIdResponse, error) {

	review, err := r.reviewUseCase.DeleteReviewById(ctx, uuid.MustParse(request.Id))
	if err != nil {
		return nil, err
	}

	response := &reviewProtobuf.DeleteReviewByIdResponse{ReviewObject: &reviewProtobuf.ReviewObject{
		Id:          review.Id.String(),
		Body:        review.Body,
		CreatedTime: review.CreatedTime.String(),
		UpdatedTime: review.UpdatedTime.String(),
	}}

	return response, nil
}

func (r *ReviewgRPC) GetCoachReviews(
	ctx context.Context,
	request *reviewProtobuf.GetCoachReviewsRequest,
) (*reviewProtobuf.GetCoachReviewsResponse, error) {
	reviews, err := r.reviewUseCase.GetCoachReviews(ctx, uuid.MustParse(request.CoachId))
	if err != nil {
		return nil, err
	}

	var reviewObjects []*reviewProtobuf.ReviewObject

	for _, review := range reviews {

		reviewObject := &reviewProtobuf.ReviewObject{
			Id:          review.Id.String(),
			UserId:      review.UserId.String(),
			Body:        review.Body,
			CreatedTime: review.CreatedTime.String(),
			UpdatedTime: review.UpdatedTime.String(),
		}

		reviewObjects = append(reviewObjects, reviewObject)
	}

	response := &reviewProtobuf.GetCoachReviewsResponse{
		ReviewObjects: reviewObjects,
	}

	return response, nil
}

func (r *ReviewgRPC) GetCoachesReviews(
	ctx context.Context,
	request *reviewProtobuf.GetCoachesReviewsRequest,
) (*reviewProtobuf.GetCoachesReviewsResponse, error) {

	var idsUUID []uuid.UUID
	for _, id := range request.CoachesIds {
		idsUUID = append(idsUUID, uuid.MustParse(id))
	}

	idsWithReviews, err := r.reviewUseCase.GetCoachesReviews(ctx, idsUUID)
	if err != nil {
		return nil, err
	}

	var coachIdWithReviewObjectArr []*reviewProtobuf.CoachIdWithReviewObject
	for key, idWithReviews := range idsWithReviews {

		coachIdWithReviewObject := &reviewProtobuf.CoachIdWithReviewObject{
			CoachId:       key.String(),
			ReviewObjects: nil,
		}

		var reviewObjects []*reviewProtobuf.ReviewObject
		for _, review := range idWithReviews {

			reviewObject := &reviewProtobuf.ReviewObject{
				Id:          review.Id.String(),
				UserId:      review.UserId.String(),
				Body:        review.Body,
				CreatedTime: review.CreatedTime.String(),
				UpdatedTime: review.UpdatedTime.String(),
			}

			reviewObjects = append(reviewObjects, reviewObject)
		}

		coachIdWithReviewObject.ReviewObjects = reviewObjects

		coachIdWithReviewObjectArr = append(coachIdWithReviewObjectArr, coachIdWithReviewObject)
	}

	response := &reviewProtobuf.GetCoachesReviewsResponse{
		CoachIdWithReviewObject: coachIdWithReviewObjectArr,
	}

	return response, nil
}
