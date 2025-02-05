package training_usecase

import (
	"Training/internal/model"
	"context"
	"fmt"
)

type TrainingRepository interface {
	Insert(context.Context, model.Training) (model.Training, error)
	UpdateTrainingsStatuses(ctx context.Context) (activeTrainings []model.Training, passedTrainings []model.Training, err error)
	GetTrainingsByDateAndCoach(ctx context.Context, date string, coachId string) ([]model.Training, error)
}

type Training struct {
	repository TrainingRepository
}

func NewTraining(repository TrainingRepository) Training {
	return Training{
		repository,
	}
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

	insertedTrainingModel, err := t.repository.Insert(ctx, trainingModel)
	if err != nil {
		return model.Training{}, err
	}

	return insertedTrainingModel, nil
}

func (t Training) UpdateRoomsList(ctx context.Context, roomMap *model.RoomMap) error {
	activeTrainings, passedTrainings, err := t.repository.UpdateTrainingsStatuses(ctx)
	if err != nil {
		return err
	}

	for _, training := range activeTrainings {
		roomMap.InitRoom(model.RoomMapKey{
			RoomId:   training.Id,
			ClientId: training.ClientId,
			CoachId:  training.CoachId,
		})
	}

	for _, training := range passedTrainings {
		roomMap.DeleteRoom(model.RoomMapKey{
			RoomId:   training.Id,
			ClientId: training.ClientId,
			CoachId:  training.CoachId,
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
