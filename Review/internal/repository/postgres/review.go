package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/DanKo-code/FitnessCenter-Review/internal/dtos"
	customErrors "github.com/DanKo-code/FitnessCenter-Review/internal/errors"
	"github.com/DanKo-code/FitnessCenter-Review/internal/models"
	"github.com/DanKo-code/FitnessCenter-Review/pkg/logger"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"time"
)

type ReviewRepository struct {
	db *sqlx.DB
}

func NewReviewRepository(db *sqlx.DB) *ReviewRepository {
	return &ReviewRepository{db: db}
}

func (r *ReviewRepository) CreateCoachReview(ctx context.Context, review *models.Review, coachId uuid.UUID) error {
	txx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}

	defer func() {
		if err != nil {
			_ = txx.Rollback()
		}
	}()

	reviewQuery := `
		INSERT INTO review (id, user_id, body, created_time, updated_time)
		VALUES (:id, :user_id, :body, :created_time, :updated_time)
	`
	_, err = txx.NamedExecContext(ctx, reviewQuery, map[string]interface{}{
		"id":           review.Id,
		"user_id":      review.UserId,
		"body":         review.Body,
		"created_time": review.CreatedTime,
		"updated_time": review.UpdatedTime,
	})
	if err != nil {
		return fmt.Errorf("failed to insert review: %w", err)
	}

	coachReviewQuery := `
		INSERT INTO coach_review (coach_id, review_id)
		VALUES (:coach_id, :review_id)
	`
	_, err = txx.NamedExecContext(ctx, coachReviewQuery, map[string]interface{}{
		"coach_id":  coachId,
		"review_id": review.Id,
	})
	if err != nil {
		return fmt.Errorf("failed to link review with coach: %w", err)
	}

	if err = txx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r *ReviewRepository) GetReviewById(ctx context.Context, id uuid.UUID) (*models.Review, error) {
	review := &models.Review{}
	err := r.db.GetContext(ctx, review, `SELECT id, user_id, body, created_time, updated_time FROM "review" WHERE id = $1`, id)
	if err != nil {
		logger.ErrorLogger.Printf("Error GetReviewById: %v", err)

		if errors.Is(err, sql.ErrNoRows) {
			return nil, customErrors.ReviewNotFound
		}

		return nil, err
	}

	return review, nil
}

func (r *ReviewRepository) UpdateReview(ctx context.Context, cmd *dtos.UpdateReviewCommand) error {
	setFields := map[string]interface{}{}

	if cmd.Body != "" {
		setFields["body"] = cmd.Body
	}
	setFields["updated_time"] = cmd.UpdatedTime

	if len(setFields) == 0 {
		logger.InfoLogger.Printf("No fields to update for review Id: %v", cmd.Id)
		return nil
	}

	query := `UPDATE "review" SET `

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

	_, err := r.db.ExecContext(ctx, query, params...)
	if err != nil {
		logger.ErrorLogger.Printf("Error UpdateReview: %v", err)
		return err
	}

	return nil
}

func (r *ReviewRepository) DeleteReviewById(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM "review" WHERE id = $1`, id)

	if err != nil {
		logger.ErrorLogger.Printf("Error DeleteReview: %v", err)
		return err
	}

	return nil
}

func (r *ReviewRepository) GetCoachReviews(ctx context.Context, coachId uuid.UUID) ([]*models.Review, error) {
	var reviews []*models.Review
	err := r.db.SelectContext(ctx, &reviews,
		`SELECT id, user_id, body, created_time, updated_time
		 FROM review
	     JOIN "coach_review" on review.id = coach_review.review_id
		 WHERE coach_review.review_id = $1`, coachId)
	if err != nil {
		return nil, err
	}

	return reviews, nil
}

func (r *ReviewRepository) GetCoachesReviews(ctx context.Context, ids []uuid.UUID) (map[uuid.UUID][]*models.Review, error) {
	coachReviews := make(map[uuid.UUID][]*models.Review)

	if len(ids) == 0 {
		return coachReviews, nil
	}

	query := `
		SELECT coach_review.coach_id, review.id, review.body, review.created_time, review.updated_time, review.user_id
		FROM "review"
		JOIN "coach_review" ON review.id = coach_review.review_id
		WHERE coach_review.coach_id = ANY($1)
	`

	type resultRow struct {
		CoachID     uuid.UUID `db:"coach_id"`
		Id          uuid.UUID `db:"id"`
		Body        string    `db:"body"`
		UserId      uuid.UUID `db:"user_id"`
		CreatedTime time.Time `db:"created_time"`
		UpdatedTime time.Time `db:"updated_time"`
	}

	var rows []resultRow

	err := r.db.SelectContext(ctx, &rows, query, pq.Array(ids))
	if err != nil {
		return nil, err
	}

	for _, row := range rows {

		review := &models.Review{
			Id:          row.Id,
			Body:        row.Body,
			UserId:      row.UserId,
			CreatedTime: row.CreatedTime,
			UpdatedTime: row.UpdatedTime,
		}

		coachReviews[row.CoachID] = append(coachReviews[row.CoachID], review)
	}

	return coachReviews, nil
}
