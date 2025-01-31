package models

import (
	"github.com/google/uuid"
	"time"
)

type RefreshSessions struct {
	Id           uuid.UUID `db:"id"`
	UserId       uuid.UUID `db:"user_id"`
	RefreshToken string    `db:"refresh_token"`
	FingerPrint  string    `db:"finger_print"`
	CreatedTime  time.Time `db:"created_time"`
	UpdatedTime  time.Time `db:"updated_time"`
}
