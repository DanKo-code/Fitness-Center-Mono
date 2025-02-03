package repository

import (
	"Training/internal/model"
	"Training/pkg/logger"
	"context"
	"github.com/jmoiron/sqlx"
	"time"
)

type Training struct {
	db sqlx.DB
}

func NewTraining(db sqlx.DB) Training {
	return Training{
		db: db,
	}
}

func (t Training) Insert(ctx context.Context, training model.Training) (model.Training, error) {

	trainingDB := convertTrainingModelToDB(training)

	query := `
				INSERT INTO training (id, time_from, time_until, status, coach_id, client_id, created_time, updated_time)
				VALUES (:id, :time_from, :time_until, :status, :coach_id, :client_id, :created_time, :updated_time)
				`

	_, err := t.db.NamedQueryContext(ctx, query, trainingDB)
	if err != nil {
		return model.Training{}, err
	}

	return training, nil
}

func (t Training) UpdateTrainingsStatuses(ctx context.Context) (activeTrainings []model.Training, passedTrainings []model.Training, err error) {

	tx, err := t.db.Beginx()
	if err != nil {
		return []model.Training{}, []model.Training{}, nil
	}

	defer func() {
		if err != nil {
			err := tx.Rollback()
			if err != nil {
				logger.Logger.Error("Error Rollback UpdateTrainingsStatuses")
			}
		}
	}()

	setActiveStatuesQuery := `
		UPDATE training
		SET status = 'active'
		WHERE status = 'booked' AND time_from <= $1
		RETURNING id, time_from, time_until, status, coach_id, client_id, created_time, updated_time`

	setPassedStatuesQuery := `
		UPDATE training
		SET status = 'passed'
		WHERE status = 'active' AND time_until <= $1
		RETURNING id, time_from, time_until, status, coach_id, client_id, created_time, updated_time`

	activeTrainingsRows, err := tx.QueryxContext(ctx, setActiveStatuesQuery, time.Now())
	if err != nil {
		return []model.Training{}, []model.Training{}, err
	}

	var activeTrainingsModels []model.Training
	for activeTrainingsRows.Next() {
		var trainingDB TrainingDB
		err = activeTrainingsRows.StructScan(&trainingDB)
		if err != nil {
			return []model.Training{}, []model.Training{}, err
		}
		trainingModel := convertTrainingDBToModel(trainingDB)
		activeTrainingsModels = append(activeTrainingsModels, trainingModel)
	}

	err = activeTrainingsRows.Close()
	if err != nil {
		return nil, nil, err
	}

	passedTrainingsRows, err := tx.QueryxContext(ctx, setPassedStatuesQuery, time.Now())
	if err != nil {
		return []model.Training{}, []model.Training{}, err
	}

	var passedTrainingsModels []model.Training
	for passedTrainingsRows.Next() {
		var trainingDB TrainingDB
		err = passedTrainingsRows.StructScan(&trainingDB)
		if err != nil {
			return []model.Training{}, []model.Training{}, err
		}
		trainingModel := convertTrainingDBToModel(trainingDB)
		passedTrainingsModels = append(passedTrainingsModels, trainingModel)
	}

	err = tx.Commit()
	if err != nil {
		return []model.Training{}, []model.Training{}, err
	}

	return activeTrainingsModels, passedTrainingsModels, nil
}

func (t Training) GetTrainingsByDate(ctx context.Context, date string) ([]model.Training, error) {

	var trainings []TrainingDB

	query := `
				SELECT id, time_from, time_until, status, coach_id, client_id, created_time, updated_time
				FROM training
				WHERE time_from::DATE = $1
			`

	err := t.db.SelectContext(ctx, &trainings, query, date)
	if err != nil {
		return nil, err
	}

	var trainingsModel []model.Training

	for _, training := range trainings {
		trainingsModel = append(trainingsModel, convertTrainingDBToModel(training))
	}

	return trainingsModel, nil
}

func convertTrainingModelToDB(trainingModel model.Training) TrainingDB {
	return TrainingDB{
		Id:          trainingModel.Id,
		TimeFrom:    trainingModel.TimeFrom,
		TimeUntil:   trainingModel.TimeUntil,
		Status:      trainingModel.Status,
		CoachId:     trainingModel.CoachId,
		ClientId:    trainingModel.ClientId,
		CreatedTime: trainingModel.CreatedTime,
		UpdatedTime: trainingModel.UpdatedTime,
	}
}

func convertTrainingDBToModel(trainingDB TrainingDB) model.Training {
	return model.Training{
		Id:          trainingDB.Id,
		TimeFrom:    trainingDB.TimeFrom,
		TimeUntil:   trainingDB.TimeUntil,
		Status:      trainingDB.Status,
		CoachId:     trainingDB.CoachId,
		ClientId:    trainingDB.ClientId,
		CreatedTime: trainingDB.CreatedTime,
		UpdatedTime: trainingDB.UpdatedTime,
	}
}
