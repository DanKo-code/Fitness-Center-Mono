package dtos

type User struct {
	Name string `json:"name" validate:"required,min=2,max=100"`
}
