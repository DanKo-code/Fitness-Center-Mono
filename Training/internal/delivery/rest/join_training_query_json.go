package rest

import "github.com/google/uuid"

type JoinTrainingQuery struct {
	CoachId uuid.UUID `json:"coach_id"`
}
