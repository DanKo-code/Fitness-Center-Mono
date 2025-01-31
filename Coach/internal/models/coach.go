package models

import (
	"github.com/google/uuid"
	"time"
)

type Coach struct {
	Id          uuid.UUID `db:"id"`
	Name        string    `db:"name"`
	Description string    `db:"description"`
	Photo       string    `db:"photo"`
	UpdatedTime time.Time `db:"updated_time"`
	CreatedTime time.Time `db:"created_time"`
}
