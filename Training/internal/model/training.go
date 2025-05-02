package model

import (
	"time"

	"github.com/google/uuid"
)

type Training struct {
	Id          uuid.UUID
	TimeFrom    time.Time
	TimeUntil   time.Time
	Status      string
	CoachId     uuid.UUID
	ClientId    uuid.UUID
	CreatedTime time.Time
	UpdatedTime time.Time
}

type Coach struct {
	Id          string
	Name        string
	Description string
	Photo       string
	UpdatedTime string
	CreatedTime string
	User        string
	Shift       string
}

type TrainingWithCoachDetails struct {
	Training Training
	Coach    Coach
}
