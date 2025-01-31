package dtos

import "github.com/google/uuid"

type CreateReviewCommand struct {
	UserId  uuid.UUID `json:"user_id"`
	Body    string    `json:"body"`
	CoachId uuid.UUID `json:"coach_id"`
}
