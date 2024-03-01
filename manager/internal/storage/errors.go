package storage

import "errors"

var (
	ErrNoSuchRequest = errors.New("no request with such id")
	ErrTooLongCrack  = errors.New("some crack is too long")
)
