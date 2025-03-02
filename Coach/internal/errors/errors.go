package errors

import "errors"

var (
	VoidCoachData      = errors.New("void coach data")
	CoachAlreadyExists = errors.New("coach already exists")
	CoachNotFound      = errors.New("coach not found")

	VoidUserData = errors.New("void userData")
)
