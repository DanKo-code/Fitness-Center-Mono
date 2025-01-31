package repository

import (
	"context"
	"github.com/DanKo-code/FitnessCenter-Order/internal/models"
	"github.com/google/uuid"
)

type OrderRepository interface {
	CreateCoachAbonement(ctx context.Context, order *models.Order) error

	GetUserOrders(ctx context.Context, userId uuid.UUID) ([]*models.Order, error)
	SetExpiredOrdersTasks(ctx context.Context) error
}
