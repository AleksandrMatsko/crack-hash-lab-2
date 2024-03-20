package config

import (
	"github.com/spf13/viper"
	"log"
	"strconv"
	"time"
)

func configureEnvs() {
	// bind envs to RabbitMQ keys
	_ = viper.BindEnv(rabbitMQConnStrKey, appEnvRabbitMQPrefix+"_CONNSTR")
	_ = viper.BindEnv(rabbitMQTaskExchangeKey, appEnvRabbitMQPrefix+"_TASK_EXCHANGE")
	_ = viper.BindEnv(rabbitMQResultExchangeKey, appEnvRabbitMQPrefix+"_RESULT_EXCHANGE")
	_ = viper.BindEnv(rabbitMQTaskQueueKey, appEnvRabbitMQPrefix+"_TASK_QUEUE")
	_ = viper.BindEnv(rabbitMQReconnectTimeoutKey, appEnvRabbitMQPrefix+"_RECONNECT_TIMEOUT")
}

func ConfigureApp() {
	viper.SetConfigName("conf")
	viper.AddConfigPath("./configs")
	viper.AddConfigPath("./worker/configs")

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("config file were not provided")
	}

	configureEnvs()
	viper.AutomaticEnv()
}

func GetRabbitMQConnStr() (string, error) {
	s := viper.GetString(rabbitMQConnStrKey)
	if s == "" {
		return "", ErrNoRabbitMQConnStr
	}
	return s, nil
}

func GetRabbitMQTaskExchange() string {
	s := viper.GetString(rabbitMQTaskExchangeKey)
	if s == "" {
		return defaultTaskExchange
	}
	return s
}

func GetRabbitMQResultExchange() string {
	s := viper.GetString(rabbitMQResultExchangeKey)
	if s == "" {
		return defaultResultExchange
	}
	return s
}

func GetRabbitMQTaskQueue() string {
	s := viper.GetString(rabbitMQTaskQueueKey)
	if s == "" {
		return defaultTaskQueue
	}
	return s
}

func GetRabbitMQReconnectTimeout() time.Duration {
	s := viper.GetString(rabbitMQReconnectTimeoutKey)
	if s == "" {
		return defaultReconnectTimeout
	}
	val, err := strconv.Atoi(s)
	if err != nil {
		return defaultReconnectTimeout
	}
	return time.Duration(val)
}
