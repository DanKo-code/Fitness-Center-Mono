package dtos

import (
	"time"

	"github.com/google/uuid"
)

type UpdateReviewCommand struct {
	Id          uuid.UUID `json:"id"`
	Body        string    `json:"body"`
	UpdatedTime time.Time `json:"updated_time"`
}
