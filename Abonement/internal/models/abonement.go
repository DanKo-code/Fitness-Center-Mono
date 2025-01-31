package models

import (
	"github.com/google/uuid"
	"time"
)

type Abonement struct {
	Id            uuid.UUID `db:"id"`
	Title         string    `db:"title"`
	Validity      string    `db:"validity"`
	VisitingTime  string    `db:"visiting_time"`
	Photo         string    `db:"photo"`
	Price         int       `db:"price"`
	UpdatedTime   time.Time `db:"updated_time"`
	CreatedTime   time.Time `db:"created_time"`
	StripePriceId string    `db:"stripe_price_id"`
}
