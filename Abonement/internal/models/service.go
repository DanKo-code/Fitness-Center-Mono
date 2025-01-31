package models

import (
	"github.com/google/uuid"
	"time"
)

type Service struct {
	Id          uuid.UUID `db:"id"`
	Title       string    `db:"title"`
	Photo       string    `db:"photo"`
	UpdatedTime time.Time `db:"updated_time"`
	CreatedTime time.Time `db:"created_time"`
}
