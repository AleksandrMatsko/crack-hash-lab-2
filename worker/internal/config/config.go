package config

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
)

func configureEnvs() {
	// bind envs to server keys
	_ = viper.BindEnv(serverHostKey, appEnvServerPrefix+"_HOST")
	_ = viper.BindEnv(serverPortKey, appEnvServerPrefix+"_PORT")

	// bind envs to manager keys
	_ = viper.BindEnv(managerHostKey, appEnvManagerPrefix+"_HOST")
	_ = viper.BindEnv(managerPortKey, appEnvManagerPrefix+"_PORT")

	// bind envs to RabbitMQ keys
	_ = viper.BindEnv(rabbitMQConnStrKey, appEnvRabbitMQPrefix+"_CONNSTR")
	_ = viper.BindEnv(rabbitMQTaskExchangeKey, appEnvRabbitMQPrefix+"_TASK_EXCHANGE")
	_ = viper.BindEnv(rabbitMQResultExchangeKey, appEnvRabbitMQPrefix+"_RESULT_EXCHANGE")
	_ = viper.BindEnv(rabbitMQTaskQueueKey, appEnvRabbitMQPrefix+"_TASK_QUEUE")
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

func GetHostPort() (string, string, error) {
	host := viper.GetString(serverHostKey)
	if host == "" {
		return "", "", ErrNoHost
	}
	port := viper.GetString(serverPortKey)
	if port == "" {
		return "", "", ErrNoPort
	}
	return host, port, nil
}

func GetManagerHostAndPort() (string, error) {
	host := viper.GetString(managerHostKey)
	if host == "" {
		return "", ErrEmptyManagerHost
	}
	port := viper.GetString(managerPortKey)
	if port == "" {
		return "", ErrEmptyManagerPort
	}
	return fmt.Sprintf("%s:%s", host, port), nil
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
