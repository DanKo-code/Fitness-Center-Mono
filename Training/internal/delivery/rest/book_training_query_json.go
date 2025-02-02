package rest

import (
	"github.com/google/uuid"
	"time"
)

type BookTrainingQuery struct {
	ClientId  uuid.UUID `json:"client_id"`
	CoachId   uuid.UUID `json:"coach_id"`
	TimeFrom  time.Time `json:"time_from"`
	TimeUntil time.Time `json:"time_until"`
}
