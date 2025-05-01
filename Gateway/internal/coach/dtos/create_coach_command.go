package dtos

type CreateCoachCommand struct {
	Name        string
	Description string
	Services    []string
	Shift       string
	Email       string
}
