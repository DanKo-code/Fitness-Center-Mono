package dtos

import (
	"time"

	"github.com/google/uuid"
)

type UpdateCoachCommand struct {
	Id          uuid.UUID `db:"id"`
	Name        string    `db:"name"`
	Description string    `db:"description"`
	Photo       string    `db:"photo"`
	UpdatedTime time.Time `db:"updated_time"`
	Shift       string    `db:"shift"`
}
