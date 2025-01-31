package postgres

import (
	"context"
	"fmt"
	"github.com/DanKo-code/FitnessCenter-Coach/internal/dtos"
	"github.com/DanKo-code/FitnessCenter-Coach/internal/models"
	"github.com/DanKo-code/FitnessCenter-Coach/internal/repository"
	"github.com/DanKo-code/FitnessCenter-Coach/pkg/logger"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

var _ repository.CoachRepository = (*CoachRepository)(nil)

type CoachRepository struct {
	db *sqlx.DB
}

func NewCoachRepository(db *sqlx.DB) *CoachRepository {
	return &CoachRepository{db: db}
}

func (coachRep *CoachRepository) CreateCoach(ctx context.Context, coach *models.Coach) (*models.Coach, error) {
	_, err := coachRep.db.NamedExecContext(ctx, `
	INSERT INTO "coach" (id, name, description, photo, created_time, updated_time)
	VALUES (:id, :name, :description, :photo, :created_time, :updated_time)`, *coach)
	if err != nil {
		logger.ErrorLogger.Printf("Error CreateCoach: %v", err)
		return nil, err
	}

	return coach, nil
}

func (coachRep *CoachRepository) GetCoachById(ctx context.Context, id uuid.UUID) (*models.Coach, error) {
	coach := &models.Coach{}
	err := coachRep.db.GetContext(ctx, coach, `SELECT id, name, description, photo, created_time, updated_time FROM "coach" WHERE id = $1`, id)
	if err != nil {
		logger.ErrorLogger.Printf("Error GetCoachById: %v", err)
		return nil, err
	}

	return coach, nil
}

func (coachRep *CoachRepository) UpdateCoach(ctx context.Context, cmd *dtos.UpdateCoachCommand) error {
	setFields := map[string]interface{}{}

	if cmd.Name != "" {
		setFields["name"] = cmd.Name
	}
	if cmd.Description != "" {
		setFields["description"] = cmd.Description
	}
	if cmd.Photo != "" {
		setFields["photo"] = cmd.Photo
	}
	setFields["updated_time"] = cmd.UpdatedTime

	if len(setFields) == 0 {
		logger.InfoLogger.Printf("No fields to update for coach Id: %v", cmd.Id)
		return nil
	}

	query := `UPDATE "coach" SET `

	var params []interface{}
	i := 1
	for field, value := range setFields {
		if i > 1 {
			query += ", "
		}

		query += fmt.Sprintf(`%s = $%d`, field, i)
		params = append(params, value)
		i++
	}
	query += fmt.Sprintf(` WHERE id = $%d`, i)
	params = append(params, cmd.Id)

	_, err := coachRep.db.ExecContext(ctx, query, params...)
	if err != nil {
		logger.ErrorLogger.Printf("Error UpdateCoach: %v", err)
		return err
	}

	return nil
}

func (coachRep *CoachRepository) DeleteCoachById(ctx context.Context, id uuid.UUID) error {
	_, err := coachRep.db.ExecContext(ctx, `
		DELETE FROM "coach"
		WHERE id = $1`, id)
	if err != nil {
		logger.ErrorLogger.Printf("Error DeleteCoach: %v", err)
		return err
	}

	return nil
}

func (coachRep *CoachRepository) GetCoaches(ctx context.Context) ([]*models.Coach, error) {
	var coaches []*models.Coach

	err := coachRep.db.SelectContext(ctx, &coaches, `SELECT id, name, description, photo, created_time, updated_time FROM "coach"`)
	if err != nil {
		logger.ErrorLogger.Printf("Error GetCoaches: %v", err)
		return nil, err
	}

	return coaches, nil
}
