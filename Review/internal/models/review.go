package models

import (
	"github.com/google/uuid"
	"time"
)

type Review struct {
	Id          uuid.UUID `db:"id"`
	UserId      uuid.UUID `db:"user_id"`
	Body        string    `db:"body"`
	UpdatedTime time.Time `db:"updated_time"`
	CreatedTime time.Time `db:"created_time"`
}
