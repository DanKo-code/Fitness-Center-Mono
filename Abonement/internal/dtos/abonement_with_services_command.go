package dtos

import (
	"github.com/DanKo-code/Fitness-Center-Abonement/internal/models"
	serviceGRPC "github.com/DanKo-code/FitnessCenter-Protobuf/gen/FitnessCenter.protobuf.service"
)

type AbonementWithServices struct {
	Abonement *models.Abonement            `db:"abonement"`
	Services  []*serviceGRPC.ServiceObject `db:"services"`
}
