package models

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID           uuid.UUID `db:"id"`
	Name         string    `db:"name"`
	Email        string    `db:"email"`
	Role         string    `db:"role"`
	PasswordHash string    `db:"password_hash"`
	Photo        string    `db:"photo"`
	CreatedTime  time.Time `db:"created_time"`
	UpdatedTime  time.Time `db:"updated_time"`
}
