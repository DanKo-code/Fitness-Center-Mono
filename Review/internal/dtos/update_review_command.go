package dtos

import (
	"github.com/google/uuid"
	"time"
)

type UpdateReviewCommand struct {
	Id          uuid.UUID `json:"id"`
	Body        string    `json:"body"`
	UpdatedTime time.Time `json:"updated_time"`
}
