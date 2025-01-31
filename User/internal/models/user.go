package models

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID           uuid.UUID `db:"id"`
	Name         string    `db:"name"`
	PasswordHash string    `db:"password_hash"`
	Email        string    `db:"email"`
	Role         string    `db:"role"`
	Photo        string    `db:"photo"`
	UpdatedTime  time.Time `db:"updated_time"`
	CreatedTime  time.Time `db:"created_time"`
}
