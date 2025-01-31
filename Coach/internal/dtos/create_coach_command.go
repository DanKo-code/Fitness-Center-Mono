package dtos

import "github.com/google/uuid"

type CreateCoachCommand struct {
	Id          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Photo       string    `json:"photo"`
}
