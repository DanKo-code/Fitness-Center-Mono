package repository

import (
	"SSO/internal/models"
	"context"
	"github.com/google/uuid"
)

type SSORepository interface {
	CreateRefreshSession(ctx context.Context, refreshSessions *models.RefreshSessions) (*models.RefreshSessions, error)
	DeleteRefreshSessionByRefreshToken(ctx context.Context, refreshToken string) error
	DeleteRefreshSessionByUserId(ctx context.Context, userId uuid.UUID) error
	GetRefreshSession(ctx context.Context, refreshToken string) (*models.RefreshSessions, error)
}
