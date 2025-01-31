package dtos

import (
	"github.com/google/uuid"
	"time"
)

type UpdateServiceCommand struct {
	Id          uuid.UUID `db:"id"`
	Title       string    `db:"title"`
	Photo       string    `db:"photo"`
	UpdatedTime time.Time `db:"updated_time"`
}
