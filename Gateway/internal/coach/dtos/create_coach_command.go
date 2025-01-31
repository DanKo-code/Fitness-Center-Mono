package dtos

type CreateCoachCommand struct {
	Name        string
	Description string
	Services    []string
}
