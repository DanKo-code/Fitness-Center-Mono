package models

import (
	"github.com/google/uuid"
	"time"
)

type Order struct {
	Id          uuid.UUID `db:"id"`
	AbonementId uuid.UUID `db:"abonement_id"`
	UserId      uuid.UUID `db:"user_id"`
	Status      string    `db:"status"`
	UpdatedTime time.Time `db:"updated_time"`
	CreatedTime time.Time `db:"created_time"`
	ExpiredTime time.Time `db:"expiration_time"`
}
