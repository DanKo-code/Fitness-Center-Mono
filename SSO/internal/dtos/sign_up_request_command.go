package dtos

type SignUpRequestCommand struct {
	Name        string
	Email       string
	Password    string
	FingerPrint string
}
