package dtos

import "github.com/google/uuid"

type CreateAbonementCommand struct {
	Id           uuid.UUID `json:"id"`
	Title        string    `json:"title"`
	Validity     string    `json:"validity"`
	VisitingTime string    `json:"visiting_time"`
	Photo        string    `json:"photo"`
	Price        int       `json:"price"`
}
