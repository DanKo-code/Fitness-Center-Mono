package training_usecase

import (
	"Training/internal/model"
	"context"
)

type TrainingRepository interface {
	Insert(context.Context, model.Training) (model.Training, error)
	UpdateTrainingsStatuses(ctx context.Context) (activeTrainings []model.Training, passedTrainings []model.Training, err error)
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
