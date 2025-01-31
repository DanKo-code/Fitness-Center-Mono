package abonement_errors

import "errors"

var (
	OnlyPhotoOptional = errors.New("только фото необязательно при создании абонемента")
)
