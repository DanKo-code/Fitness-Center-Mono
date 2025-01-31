package dtos

import "github.com/google/uuid"

type CreateCoachServicesCommand struct {
	CoachId     uuid.UUID   `db:"coach_id"`
	ServicesIds []uuid.UUID `db:"services_ids"`
}
