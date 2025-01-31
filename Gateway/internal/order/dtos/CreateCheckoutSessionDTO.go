package dtos

type CreateCheckoutSessionDTO struct {
	UserId        string
	ClientId      string
	AbonementId   string
	StripePriceId string
}
