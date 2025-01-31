package dtos

type UpdateAbonementCommand struct {
	Id             string
	Title          string
	ValidityPeriod string
	VisitingTime   string
	Price          int64
	Services       []string
}
