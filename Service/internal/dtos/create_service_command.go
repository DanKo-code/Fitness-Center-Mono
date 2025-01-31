package dtos

import "github.com/google/uuid"

type CreateServiceCommand struct {
	Id    uuid.UUID `json:"id"`
	Title string    `db:"title"`
	Photo string    `db:"photo"`
}
