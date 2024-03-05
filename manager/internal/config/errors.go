package config

import "errors"

var (
	ErrNoHost = errors.New("no host to bind")
	ErrNoPort = errors.New("no port to bind")

	ErrBadWorkersListFormat = errors.New("bad workers.list format")

	ErrNoMongoConnStr = errors.New("mongo connection string not provided")

	ErrNoRabbitMQConnStr = errors.New("rabbitmq connection string not provided")
)
