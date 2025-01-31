package postgres

import (
	"context"
	"fmt"
	"github.com/DanKo-code/Fitness-Center-Abonement/internal/dtos"
	"github.com/DanKo-code/Fitness-Center-Abonement/internal/models"
	"github.com/DanKo-code/Fitness-Center-Abonement/internal/repository"
	"github.com/DanKo-code/Fitness-Center-Abonement/pkg/logger"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

var _ repository.AbonementRepository = (*AbonementRepository)(nil)

type AbonementRepository struct {
	db *sqlx.DB
}

func NewAbonementRepository(db *sqlx.DB) *AbonementRepository {
	return &AbonementRepository{db: db}
}

func (abonementRep *AbonementRepository) CreateAbonement(ctx context.Context, abonement *models.Abonement) (*models.Abonement, error) {
	_, err := abonementRep.db.NamedExecContext(ctx, `
	INSERT INTO "abonement" (id, title, validity,visiting_time, photo, price, created_time, updated_time, stripe_price_id)
	VALUES (:id, :title, :validity, :visiting_time, :photo, :price, :created_time, :updated_time, :stripe_price_id)`, *abonement)
	if err != nil {
		logger.ErrorLogger.Printf("Error CreateAbonement: %v", err)
		return nil, err
	}

	return abonement, nil
}

func (abonementRep *AbonementRepository) GetAbonementById(ctx context.Context, id uuid.UUID) (*models.Abonement, error) {
	abonement := &models.Abonement{}
	err := abonementRep.db.GetContext(ctx, abonement, `
		SELECT id, title, validity,visiting_time, photo, price, created_time, updated_time, stripe_price_id
		FROM "abonement"
		WHERE id = $1`, id)
	if err != nil {
		logger.ErrorLogger.Printf("Error GetAbonementById: %v", err)
		return nil, err
	}

	return abonement, nil
}

func (abonementRep *AbonementRepository) UpdateAbonement(ctx context.Context, cmd *dtos.UpdateAbonementCommand) error {
	setFields := map[string]interface{}{}

	if cmd.Title != "" {
		setFields["title"] = cmd.Title
	}
	if cmd.Validity != "" {
		setFields["validity"] = cmd.Validity
	}
	if cmd.VisitingTime != "" {
		setFields["visiting_time"] = cmd.VisitingTime
	}
	if cmd.Photo != "" {
		setFields["photo"] = cmd.Photo
	}
	if cmd.Price != 0 {
		setFields["price"] = cmd.Price
	}
	if cmd.StripePriceId != "" {
		setFields["stripe_price_id"] = cmd.StripePriceId
	}
	setFields["updated_time"] = cmd.UpdatedTime

	if len(setFields) == 0 {
		logger.InfoLogger.Printf("No fields to update for abonement Id: %v", cmd.Id)
		return nil
	}

	query := `UPDATE "abonement" SET `

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

	_, err := abonementRep.db.ExecContext(ctx, query, params...)
	if err != nil {
		logger.ErrorLogger.Printf("Error UpdateAbonement: %v", err)
		return err
	}

	return nil
}

func (abonementRep *AbonementRepository) DeleteAbonementById(ctx context.Context, id uuid.UUID) error {
	_, err := abonementRep.db.ExecContext(ctx, `
		DELETE FROM "abonement"
		WHERE id = $1`, id)
	if err != nil {
		logger.ErrorLogger.Printf("Error DeleteAbonement: %v", err)
		return err
	}

	return nil
}

func (abonementRep *AbonementRepository) GetAbonementes(ctx context.Context) ([]*models.Abonement, error) {
	var abonementes []*models.Abonement

	err := abonementRep.db.SelectContext(ctx, &abonementes, `SELECT id, title, validity,visiting_time, photo, price, created_time, updated_time, stripe_price_id  FROM "abonement"`)
	if err != nil {
		logger.ErrorLogger.Printf("Error GetAbonementes: %v", err)
		return nil, err
	}

	return abonementes, nil
}

func (abonementRep *AbonementRepository) GetAbonementsByIds(ctx context.Context, ids []uuid.UUID) ([]*models.Abonement, error) {
	query := `SELECT id, title, validity, visiting_time, photo, price, created_time, updated_time, stripe_price_id 
			  FROM "abonement"
			  WHERE id IN (?)`

	query, args, err := sqlx.In(query, ids)
	if err != nil {
		return nil, fmt.Errorf("failed to bind ids: %w", err)
	}

	query = abonementRep.db.Rebind(query)

	var abonements []*models.Abonement

	err = abonementRep.db.SelectContext(ctx, &abonements, query, args...)
	if err != nil {
		return nil, err
	}

	if len(abonements) != len(ids) {

		idAbonementMap := make(map[uuid.UUID]*models.Abonement, len(abonements))

		for _, abonement := range abonements {
			idAbonementMap[abonement.Id] = abonement
		}

		var res []*models.Abonement

		for _, id := range ids {
			res = append(res, idAbonementMap[id])
		}

		return res, nil
	}

	return abonements, nil
}
