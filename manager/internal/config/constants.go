package config

// environment variables prefixes, used for app configuration
const (
	appEnvPrefix        = "MANAGER"
	appEnvServerPrefix  = appEnvPrefix + "_SERVER"
	appEnvWorkersPrefix = appEnvPrefix + "_WORKERS"
)

type RequestStatus string

// request statuses
const (
	Ready      RequestStatus = "READY"
	InProgress RequestStatus = "IN_PROGRESS"
	Error      RequestStatus = "ERROR"
)
