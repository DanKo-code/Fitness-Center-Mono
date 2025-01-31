package coach_errors

import "errors"

var (
	OnlyPhotoOptional = errors.New("при создании тренера необязательна только фотография")
)
