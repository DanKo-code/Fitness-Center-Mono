package dtos

type RefreshRequestCommand struct {
	FingerPrint  string `json:"finger_print"`
	RefreshToken string `json:"refresh_token"`
}
