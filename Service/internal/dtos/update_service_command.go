package dtos

import (
	"time"

	"github.com/google/uuid"
)

type UpdateServiceCommand struct {
	Id          uuid.UUID `db:"id"`
	Title       string    `db:"title"`
	Photo       string    `db:"photo"`
	UpdatedTime time.Time `db:"updated_time"`
}
