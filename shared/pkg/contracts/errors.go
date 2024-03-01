package contracts

import "errors"

var (
	ErrEmptyAlphabet     = errors.New("empty alphabet")
	ErrNegativeMaxLength = errors.New("negative max length")
	ErrEmptyHashToCrack  = errors.New("empty hash to crack")
)
