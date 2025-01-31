package postgres

import (
	"SSO/internal/models"
	"SSO/internal/ssoErrors"
	"SSO/pkg/logger"
	"context"
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type SSORepository struct {
	db *sqlx.DB
}

func NewSSORepository(db *sqlx.DB) *SSORepository {
	return &SSORepository{db: db}
}

func (ssoRep *SSORepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	user := &models.User{}
	err := ssoRep.db.GetContext(ctx, user, `SELECT id, name, email, role, password_hash, photo, created_time, updated_time FROM "user" WHERE email = $1`, email)
	if err != nil {
		logger.ErrorLogger.Printf("Error GetUserByEmail: %v", err)

		if errors.Is(err, sql.ErrNoRows) {
			return nil, ssoErrors.UserNotFound
		}

		return nil, err
	}

	return user, nil
}

func (ssoRep *SSORepository) CreateRefreshSession(ctx context.Context, refreshSessions *models.RefreshSessions) (*models.RefreshSessions, error) {
	_, err := ssoRep.db.NamedExecContext(ctx, `
		INSERT INTO refresh_sessions (id, user_id, refresh_token, finger_print, created_time, updated_time)
		VALUES (:id, :user_id, :refresh_token, :finger_print, :created_time, :updated_time)`, *refreshSessions)
	if err != nil {
		logger.ErrorLogger.Printf("Error CreateRefreshSession: %v", err)
		return nil, err
	}

	return refreshSessions, nil
}

func (ssoRep *SSORepository) DeleteRefreshSessionByRefreshToken(ctx context.Context, refreshToken string) error {

	params := map[string]interface{}{
		"RefreshToken": refreshToken,
	}

	_, err := ssoRep.db.NamedExecContext(ctx, `DELETE FROM refresh_sessions WHERE refresh_token = :RefreshToken`, params)
	if err != nil {
		logger.ErrorLogger.Printf("Error DeleteRefreshSessionByRefreshToken: %v", err)
		return err
	}

	return nil
}

func (ssoRep *SSORepository) DeleteRefreshSessionByUserId(ctx context.Context, userId uuid.UUID) error {

	res, err := ssoRep.db.ExecContext(ctx, "DELETE FROM refresh_sessions WHERE user_id = $1", userId)
	if err != nil {
		logger.ErrorLogger.Printf("Error DeleteRefreshSessionByUserId: %v", err)
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		logger.ErrorLogger.Printf("Error RowsAffected: %v", err)
		return err
	}

	if affected == 0 {
		logger.ErrorLogger.Printf("Error refresh session not found: %v", err)
		return ssoErrors.RefreshSessionNotFound
	}

	return nil
}

func (ssoRep *SSORepository) GetRefreshSession(ctx context.Context, refreshToken string) (*models.RefreshSessions, error) {
	refreshSession := &models.RefreshSessions{}
	err := ssoRep.db.GetContext(ctx, refreshSession, `SELECT finger_print FROM refresh_sessions WHERE refresh_token = $1`, refreshToken)
	if err != nil {
		logger.ErrorLogger.Printf("Error GetRefreshSession: %v", err)
		return nil, err
	}

	return refreshSession, nil
}
