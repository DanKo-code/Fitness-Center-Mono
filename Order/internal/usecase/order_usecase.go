package usecase

import (
	"context"
	"github.com/DanKo-code/FitnessCenter-Order/internal/dtos"
	"github.com/DanKo-code/FitnessCenter-Order/internal/models"
	orderGRPC "github.com/DanKo-code/FitnessCenter-Protobuf/gen/FitnessCenter.protobuf.order"
	"github.com/google/uuid"
)

type OrderUseCase interface {
	CreateCoachAbonement(ctx context.Context, cmd *dtos.CreateOrderCommand) (*models.Order, error)
	GetUserOrders(ctx context.Context, userId uuid.UUID) (*orderGRPC.GetUserOrdersResponse, error)
	SetExpiredOrdersTasks(ctx context.Context) error
}
