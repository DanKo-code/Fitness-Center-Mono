package errors

import "errors"

var (
	ErrTrainingNotFound     = errors.New("training not found")
	ErrNotValidTrainingData = errors.New("not valid training data")
)
