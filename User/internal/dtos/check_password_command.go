package dtos

import "github.com/google/uuid"

type CheckPasswordCommand struct {
	UserId   uuid.UUID `json:"user_id"`
	Password string    `json:"password"`
}
