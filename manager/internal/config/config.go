package config

import (
	"github.com/spf13/viper"
	"log"
	"strconv"
	"time"
)

// configureEnvs bind viper keys to specified envs
func configureEnvs() {
	// bind envs to server keys
	_ = viper.BindEnv(serverHostKey, appEnvServerPrefix+"_HOST")
	_ = viper.BindEnv(serverPortKey, appEnvServerPrefix+"_PORT")

	// bind envs to workers keys
	_ = viper.BindEnv(workersTaskNumParts, appEnvWorkersPrefix+"_TASK_NUM_PARTS")

	// bind envs to MongoDB keys
	_ = viper.BindEnv(mongoConnStrKey, appEnvMongoPrefix+"_CONNSTR")
	_ = viper.BindEnv(mongoDBNameKey, appEnvMongoPrefix+"_DBNAME")

	// bind envs to RabbitMQ keys
	_ = viper.BindEnv(rabbitMQConnStrKey, appEnvRabbitMQPrefix+"_CONNSTR")
	_ = viper.BindEnv(rabbitMQTaskExchangeKey, appEnvRabbitMQPrefix+"_TASK_EXCHANGE")
	_ = viper.BindEnv(rabbitMQResultExchangeKey, appEnvRabbitMQPrefix+"_RESULT_EXCHANGE")
	_ = viper.BindEnv(rabbitMQResultQueueKey, appEnvRabbitMQPrefix+"_RESULT_QUEUE")
	_ = viper.BindEnv(rabbitMQReconnectTimeoutKey, appEnvRabbitMQPrefix+"_RECONNECT_TIMEOUT")
}

func ConfigureApp() {
	viper.SetConfigName("conf")
	viper.AddConfigPath("./configs")
	viper.AddConfigPath("./manager/configs")

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

func GetTaskNumParts() uint64 {
	s := viper.GetString(workersTaskNumParts)
	val, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return defaultTaskNumParts
	}
	return val
}

func GetMongoConnStr() (string, error) {
	s := viper.GetString(mongoConnStrKey)
	if s == "" {
		return s, ErrNoMongoConnStr
	}
	return s, nil
}

func GetMongoDBName() string {
	s := viper.GetString(mongoDBNameKey)
	if s == "" {
		return defaultDBName
	}
	return s
}

func GetMongoDBCollectionName() string {
	return defaultCollectionName
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

func GetRabbitMQResultQueue() string {
	s := viper.GetString(rabbitMQResultQueueKey)
	if s == "" {
		return defaultResultQueue
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

var toSendChan chan<- []byte

func SetToSendChan(ch chan<- []byte) {
	toSendChan = ch
}

func GetToSendChan() chan<- []byte {
	return toSendChan
}
