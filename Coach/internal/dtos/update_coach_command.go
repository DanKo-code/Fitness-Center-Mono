package dtos

import (
	"github.com/google/uuid"
	"time"
)

type UpdateCoachCommand struct {
	Id          uuid.UUID `db:"id"`
	Name        string    `db:"name"`
	Description string    `db:"description"`
	Photo       string    `db:"photo"`
	UpdatedTime time.Time `db:"updated_time"`
}
