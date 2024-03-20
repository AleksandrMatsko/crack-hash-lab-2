package config

import "time"

const (
	appEnvPrefix         = "WORKER"
	appEnvServerPrefix   = appEnvPrefix + "_SERVER"
	appEnvRabbitMQPrefix = appEnvPrefix + "_RABBITMQ"
)

const (
	// server keys
	serverHostKey = "server.host"
	serverPortKey = "server.port"

	// RabbitMQ keys
	rabbitMQConnStrKey          = "rabbitmq.connStr"
	rabbitMQTaskExchangeKey     = "rabbitmq.taskExchange"
	rabbitMQResultExchangeKey   = "rabbitmq.resultExchange"
	rabbitMQTaskQueueKey        = "rabbitmq.taskQueue"
	rabbitMQReconnectTimeoutKey = "rabbitmq.reconnectTimeout"
)

const (
	// default values for RabbitMQ
	defaultTaskExchange     = "tasks"
	defaultResultExchange   = "results"
	defaultTaskQueue        = "task_queue"
	defaultReconnectTimeout = time.Second * 10
)
