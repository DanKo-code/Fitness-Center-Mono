package dtos

import (
	"time"

	"github.com/google/uuid"
)

type UpdateUserCommand struct {
	Id          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	Role        string    `json:"role"`
	Photo       string    `json:"photo"`
	UpdatedTime time.Time `db:"updated_time"`
}
