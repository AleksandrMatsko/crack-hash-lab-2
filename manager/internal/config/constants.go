package config

// environment variables prefixes, used for app configuration
const (
	appEnvPrefix        = "MANAGER"
	appEnvServerPrefix  = appEnvPrefix + "_SERVER"
	appEnvWorkersPrefix = appEnvPrefix + "_WORKERS"
	appEnvMongoPrefix   = appEnvPrefix + "_MONGO"
)

type RequestStatus string

// request statuses
const (
	Ready      RequestStatus = "READY"
	InProgress RequestStatus = "IN_PROGRESS"
	Error      RequestStatus = "ERROR"
)

const (
	serverHostKey      = "server.host"
	serverPortKey      = "server.port"
	workersListKey     = "workers.list"
	workersTaskSizeKey = "workers.taskSize"
	mongoConnStrKey    = "mongo.connStr"
	mongoDbNameKey     = "mongo.dbname"
)
