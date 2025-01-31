package dtos

type CreateAbonementCommand struct {
	Title          string
	ValidityPeriod string
	VisitingTime   string
	Price          int64
	Services       []string
}
