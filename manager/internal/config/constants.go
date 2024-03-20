package config

import "time"

// environment variables prefixes, used for app configuration
const (
	appEnvPrefix         = "MANAGER"
	appEnvServerPrefix   = appEnvPrefix + "_SERVER"
	appEnvWorkersPrefix  = appEnvPrefix + "_WORKERS"
	appEnvMongoPrefix    = appEnvPrefix + "_MONGO"
	appEnvRabbitMQPrefix = appEnvPrefix + "_RABBITMQ"
)

type RequestStatus string

// request statuses
const (
	Ready      RequestStatus = "READY"
	InProgress RequestStatus = "IN_PROGRESS"
	Error      RequestStatus = "ERROR"
)

const (
	// server keys
	serverHostKey = "server.host"
	serverPortKey = "server.port"

	// workers keys
	workersTaskNumParts = "workers.taskNumParts"

	// MongoDB keys
	mongoConnStrKey = "mongo.connStr"
	mongoDBNameKey  = "mongo.dbname"

	// RabbitMQ keys
	rabbitMQConnStrKey          = "rabbitmq.connStr"
	rabbitMQTaskExchangeKey     = "rabbitmq.taskExchange"
	rabbitMQResultExchangeKey   = "rabbitmq.resultExchange"
	rabbitMQResultQueueKey      = "rabbitmq.resultQueue"
	rabbitMQReconnectTimeoutKey = "rabbitmq.reconnectTimeout"
)

const (
	// default values for workers
	defaultTaskNumParts = 10

	// default values for MongoDB
	defaultDBName         = "CrackHash"
	defaultCollectionName = "Requests"

	// default values for RabbitMQ
	defaultTaskExchange     = "tasks"
	defaultResultExchange   = "results"
	defaultResultQueue      = "res_queue"
	defaultReconnectTimeout = time.Second * 10
)
