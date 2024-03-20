package config

import "time"

const (
	appEnvPrefix         = "WORKER"
	appEnvRabbitMQPrefix = appEnvPrefix + "_RABBITMQ"
)

const (
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
