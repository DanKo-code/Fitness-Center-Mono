package dtos

type SignInRequest struct {
	Email       string `json:"email" validate:"required,email"`
	Password    string `json:"password" validate:"required,min=8,max=30"`
	FingerPrint string `json:"finger_print"`
}

type SignInRequestWithoutFingerprint struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=30"`
}

type SignInResponse struct {
	AccessToken           string `json:"accessToken"`
	AccessTokenExpiration int    `json:"accessTokenExpiration"`
	User                  User   `json:"user"`
}
