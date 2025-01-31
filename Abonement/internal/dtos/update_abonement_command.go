package dtos

import (
	"github.com/google/uuid"
	"time"
)

type UpdateAbonementCommand struct {
	Id            uuid.UUID `json:"id"`
	Title         string    `json:"title"`
	Validity      string    `json:"validity"`
	VisitingTime  string    `json:"visiting_time"`
	Photo         string    `json:"photo"`
	UpdatedTime   time.Time `json:"updated_time"`
	Price         int       `json:"price"`
	StripePriceId string    `json:"stripe_price_id"`
}
