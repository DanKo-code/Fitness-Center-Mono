package repository

import "github.com/google/uuid"

type RoomKeyDB struct {
	RoomId   uuid.UUID `db:"room_id"`
	ClientId uuid.UUID `db:"client_id"`
	CoachId  uuid.UUID `db:"coach_id"`
}
