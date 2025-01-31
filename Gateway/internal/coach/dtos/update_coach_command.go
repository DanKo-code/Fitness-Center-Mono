package dtos

type UpdateCoachCommand struct {
	Id          string
	Name        string
	Description string
	Services    []string
}
