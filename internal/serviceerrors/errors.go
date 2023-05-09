package serviceerrors

import (
	"errors"
)

var (
	ErrBlankField   = errors.New("this field cannot be blank")
	ErrLongField    = errors.New("this field cannot be more than 100 characters long")
	ErrExpiresField = errors.New("his field must equal 1, 7 or 365")
)
