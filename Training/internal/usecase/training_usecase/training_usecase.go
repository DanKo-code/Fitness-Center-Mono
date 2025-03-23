package training_usecase

import (
	errorsx "Training/internal/errors"
	"Training/internal/model"
	"context"
	"database/sql"
	"errors"
	"fmt"
	coachGRPC "github.com/DanKo-code/FitnessCenter-Protobuf/gen/FitnessCenter.protobuf.coach"
	userGRPC "github.com/DanKo-code/FitnessCenter-Protobuf/gen/FitnessCenter.protobuf.user"
	"github.com/google/uuid"
	"log/slog"
	"time"
)

type TrainingRepository interface {
	Insert(context.Context, model.Training) (model.Training, error)
	UpdateTrainingsStatuses(ctx context.Context) (activeTrainings []model.Training, passedTrainings []model.Training, err error)
	GetTrainingsByDateAndCoach(ctx context.Context, date string, coachId string) ([]model.Training, error)
	GetTrainingByTime(ctx context.Context, timeFrom, timeUntil time.Time) (model.Training, error)
	GetAvailableCoaches(ctx context.Context, training model.Training) ([]string, error)
	GetById(ctx context.Context, id uuid.UUID) (model.Training, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type Training struct {
	repository  TrainingRepository
	coachClient coachGRPC.CoachClient
	userClient  userGRPC.UserClient
}

func NewTraining(
	repository TrainingRepository,
	coachClient coachGRPC.CoachClient,
	userClient userGRPC.UserClient,
) Training {
	return Training{
		repository:  repository,
		coachClient: coachClient,
		userClient:  userClient,
	}
}

func (t Training) IsValidParticipant(ctx context.Context, roomId, clientId, coachId uuid.UUID, role string) error {
	training, err := t.repository.GetById(ctx, roomId)
	if err != nil {
		return nil
	}

	if role != "coach" {
		if training.ClientId != clientId {
			slog.Info("[Training Service] [IsValidParticipant]: invalid client id")
			return fmt.Errorf("[Training Service] [IsValidParticipant]: %w", errorsx.ErrNotValidTrainingData)
		}
		if training.Status != "active" {
			slog.Info("[Training Service] [IsValidParticipant]: status not active")
			return fmt.Errorf("[Training Service] [IsValidParticipant]: %w", errorsx.ErrNotValidTrainingData)
		}
		if training.CoachId != coachId {
			slog.Info("[Training Service] [IsValidParticipant]: invalid coach id")
			return fmt.Errorf("[Training Service] [IsValidParticipant]: %w", errorsx.ErrNotValidTrainingData)
		}
	} else {
		if training.Status != "active" {
			slog.Info("[Training Service] [IsValidParticipant]: status not active")
			return fmt.Errorf("[Training Service] [IsValidParticipant]: %w", errorsx.ErrNotValidTrainingData)
		}
		if training.CoachId != coachId {
			slog.Info("[Training Service] [IsValidParticipant]: invalid coach id")
			return fmt.Errorf("[Training Service] [IsValidParticipant]: %w", errorsx.ErrNotValidTrainingData)
		}
		coach, err := t.coachClient.GetCoachById(ctx, &coachGRPC.GetCoachByIdRequest{Id: coachId.String()})
		if err != nil {
			return fmt.Errorf("[Training Service] [IsValidParticipant]: %w", err)
		}
		if coach.CoachObject.User != clientId.String() {
			slog.Info("[Training Service] [IsValidParticipant]: invalid coach id in user id")
			return fmt.Errorf("[Training Service] [IsValidParticipant]: %w", errorsx.ErrNotValidTrainingData)
		}
	}

	return nil
}

func (t Training) Insert(ctx context.Context, trainingModel model.Training) (model.Training, error) {

	trainings, err := t.repository.GetTrainingsByDateAndCoach(ctx, trainingModel.TimeUntil.Format("2006-01-02"), trainingModel.CoachId.String())
	if err != nil {
		return model.Training{}, err
	}

	trainingsPerDayCount := 0
	for _, training := range trainings {
		if training.ClientId == trainingModel.ClientId {
			trainingsPerDayCount++
		}
	}

	if trainingsPerDayCount == 2 {
		return model.Training{}, fmt.Errorf("превышен лимит дневных тренировок: 2")
	}

	_, err = t.repository.GetTrainingByTime(ctx, trainingModel.TimeFrom, trainingModel.TimeUntil)
	if err == nil || !errors.Is(err, sql.ErrNoRows) {
		return model.Training{}, fmt.Errorf("У вас уже забронирована тренировка в это время")
	}

	_, err = t.repository.GetAvailableCoaches(ctx, trainingModel)
	if err != nil {
		return model.Training{}, err
	}

	insertedTrainingModel, err := t.repository.Insert(ctx, trainingModel)
	if err != nil {
		return model.Training{}, err
	}

	return insertedTrainingModel, nil
}

func (t Training) Delete(ctx context.Context, id uuid.UUID) error {
	err := t.repository.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (t Training) UpdateRoomsList(ctx context.Context, roomMap *model.RoomMap) error {
	activeTrainings, passedTrainings, err := t.repository.UpdateTrainingsStatuses(ctx)
	if err != nil {
		return err
	}

	for _, training := range activeTrainings {
		roomMap.InitRoom(model.RoomMapKey{
			RoomId:   training.Id,
			ClientId: uuid.Nil,
			CoachId:  uuid.Nil,
		})
	}

	for _, training := range passedTrainings {
		roomMap.DeleteRoom(model.RoomMapKey{
			RoomId:   training.Id,
			ClientId: uuid.Nil,
			CoachId:  uuid.Nil,
		})
	}

	return nil
}

func (t Training) GetTrainingsByDateAndCoach(ctx context.Context, date string, coachId string) ([]model.Training, error) {
	trainings, err := t.repository.GetTrainingsByDateAndCoach(ctx, date, coachId)
	if err != nil {
		return nil, err
	}

	return trainings, nil
}
