package dtos

import (
	"github.com/google/uuid"
)

type CreateUserCommand struct {
	ID       uuid.UUID `db:"id"`
	Name     string    `db:"name"`
	Email    string    `db:"email"`
	Role     string    `db:"role"`
	Photo    string    `db:"photo"`
	Password string    `db:"password"`
}
