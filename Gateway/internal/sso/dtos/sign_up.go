package dtos

type SignUpRequest struct {
	Name        string `json:"name" validate:"required,min=2,max=100"`
	Email       string `json:"email" validate:"required,email"`
	Password    string `json:"password" validate:"required,min=8,max=30"`
	FingerPrint string `json:"finger_print"`
}

type SignUpRequestWithOutFingerPrint struct {
	Name     string `json:"name" validate:"required,min=2,max=100"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=30"`
}

type SignUpResponse struct {
	AccessToken           string `json:"accessToken"`
	AccessTokenExpiration int    `json:"accessTokenExpiration"`
	User                  User   `json:"user"`
}
