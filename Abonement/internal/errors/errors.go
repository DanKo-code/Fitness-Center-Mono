package errors

import "errors"

var (
	VoidAbonementData      = errors.New("void abonement data")
	AbonementAlreadyExists = errors.New("abonement already exists")
	AbonementNotFound      = errors.New("abonement not found")
)
