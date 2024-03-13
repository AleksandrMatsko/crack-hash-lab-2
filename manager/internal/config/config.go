package config

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
	"strconv"
	"strings"
	"time"
)

// configureEnvs bind viper keys to specified envs
func configureEnvs() {
	// bind envs to server keys
	_ = viper.BindEnv(serverHostKey, appEnvServerPrefix+"_HOST")
	_ = viper.BindEnv(serverPortKey, appEnvServerPrefix+"_PORT")

	// bind envs to workers keys
	_ = viper.BindEnv(workersListKey, appEnvWorkersPrefix+"_LIST")
	_ = viper.BindEnv(workersTaskSizeKey, appEnvWorkersPrefix+"_TASK_SIZE")
	_ = viper.BindEnv(workersCountKey, appEnvWorkersPrefix+"_COUNT")

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

func GetWorkers() ([]string, error) {
	s := viper.GetString(workersListKey)
	if s == "" {
		return make([]string, 0), nil
	}
	splited := strings.Split(s, ":")
	if len(splited)%2 != 0 {
		return nil, ErrBadWorkersListFormat
	}
	workers := make([]string, 0)
	for i := 0; i < len(splited); i += 2 {
		if splited[i] == "" {
			log.Printf("bad worker host: '%s'", splited[i])
			continue
		}
		if splited[i+1] == "" {
			log.Printf("bad worker port: '%s'", splited[i+1])
			continue
		}
		workers = append(workers, fmt.Sprintf("%s:%s", splited[i], splited[i+1]))
	}
	return workers, nil
}

func GetTaskSize() (uint64, error) {
	s := viper.GetString(workersTaskSizeKey)
	val, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0, err
	}
	return val, nil
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
