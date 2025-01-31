package dtos

import "github.com/google/uuid"

type CreateAbonementServicesCommand struct {
	AbonementId uuid.UUID   `db:"abonement_id"`
	ServicesIds []uuid.UUID `db:"services_ids"`
}
