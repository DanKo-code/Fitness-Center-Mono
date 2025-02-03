package rest

import (
	"github.com/google/uuid"
)

type BookTrainingQuery struct {
	CoachId   uuid.UUID `json:"coach_id"`
	TimeFrom  string    `json:"time_from"`
	TimeUntil string    `json:"time_until"`
	Date      string    `json:"date"`
}
