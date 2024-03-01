package config

import "errors"

var (
	ErrEmptyManagerHost = errors.New("empty manager.host")
	ErrEmptyManagerPort = errors.New("empty manager.port")
)
