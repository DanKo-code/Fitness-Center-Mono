package postgres

import (
	"context"
	"fmt"
	"github.com/DanKo-code/FitnessCenter-Order/internal/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"time"
)

type OrderRepository struct {
	db *sqlx.DB
}

func NewOrderRepository(db *sqlx.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (o *OrderRepository) CreateCoachAbonement(ctx context.Context, order *models.Order) error {
	orderQuery := `
		INSERT INTO "order"(id, abonement_id, user_id, status, created_time, updated_time, expiration_time)
		VALUES (:id, :abonement_id, :user_id, :status, :created_time, :updated_time, :expiration_time)
	`
	_, err := o.db.NamedExecContext(ctx, orderQuery, map[string]interface{}{
		"id":              order.Id,
		"abonement_id":    order.AbonementId,
		"user_id":         order.UserId,
		"status":          order.Status,
		"created_time":    order.CreatedTime,
		"updated_time":    order.UpdatedTime,
		"expiration_time": order.ExpiredTime,
	})
	if err != nil {
		return fmt.Errorf("failed to insert order: %w", err)
	}

	return nil
}

func (o *OrderRepository) GetUserOrders(ctx context.Context, userId uuid.UUID) ([]*models.Order, error) {
	var orders []*models.Order
	err := o.db.SelectContext(ctx, &orders,
		`SELECT id, abonement_id, user_id, status, created_time, updated_time
		 FROM "order"
		 WHERE user_id = $1`, userId)
	if err != nil {
		return nil, err
	}

	return orders, nil
}

func (o *OrderRepository) SetExpiredOrdersTasks(ctx context.Context) error {
	q := `UPDATE "order" 
		  SET status = 'Expired' 
		  WHERE "order".expiration_time <= $1 AND status = 'Valid'`

	_, err := o.db.ExecContext(ctx, q, time.Now())
	if err != nil {
		return err
	}

	return nil
}
