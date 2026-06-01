package models

import (
	"errors"
)


var(
	ErrDuplicatePhone = errors.New("models:phone number already exists")
	ErrInvalidCredentials = errors.New("models: invalid credentials")

)
