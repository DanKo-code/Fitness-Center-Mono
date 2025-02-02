package model

import (
	"github.com/google/uuid"
	"time"
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
