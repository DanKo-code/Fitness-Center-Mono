package dtos

import "github.com/google/uuid"

type CreateOrderCommand struct {
	UserId      uuid.UUID `json:"user_id"`
	AbonementId uuid.UUID `json:"abonement_id"`
}
