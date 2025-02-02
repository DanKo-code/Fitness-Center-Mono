package repository

import (
	"github.com/google/uuid"
	"time"
)

type TrainingDB struct {
	Id          uuid.UUID `db:"id"`
	TimeFrom    time.Time `db:"time_from"`
	TimeUntil   time.Time `db:"time_until"`
	Status      string    `db:"status"`
	CoachId     uuid.UUID `db:"coach_id"`
	ClientId    uuid.UUID `db:"client_id"`
	CreatedTime time.Time `db:"created_time"`
	UpdatedTime time.Time `db:"updated_time"`
}
