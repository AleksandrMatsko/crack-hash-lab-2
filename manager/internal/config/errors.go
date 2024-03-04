package config

import "errors"

var (
	ErrBadWorkersListFormat = errors.New("bad workers.list format")
	ErrNoMongoConnStr       = errors.New("mongo connection string not provided")
)
