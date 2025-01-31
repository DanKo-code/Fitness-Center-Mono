package dtos

type User struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	Role        string `json:"role"`
	Photo       string `json:"photo"`
	CreatedTime string `json:"created_time"`
	UpdatedTime string `json:"updated_time"`
}
