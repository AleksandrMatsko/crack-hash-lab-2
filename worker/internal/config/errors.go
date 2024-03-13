package config

import "errors"

var (
	ErrNoHost = errors.New("no host to bind")
	ErrNoPort = errors.New("no port to bind")

	ErrEmptyManagerHost = errors.New("empty manager.host")
	ErrEmptyManagerPort = errors.New("empty manager.port")

	ErrNoRabbitMQConnStr = errors.New("rabbitmq connection string not provided")
)
