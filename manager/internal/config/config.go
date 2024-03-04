package config

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
	"strconv"
	"strings"
)

// configureEnvs bind viper keys to specified envs
func configureEnvs() {
	_ = viper.BindEnv(serverHostKey, appEnvServerPrefix+"_HOST")
	_ = viper.BindEnv(serverPortKey, appEnvServerPrefix+"_PORT")
	_ = viper.BindEnv(workersListKey, appEnvWorkersPrefix+"_LIST")
	_ = viper.BindEnv(workersTaskSizeKey, appEnvWorkersPrefix+"_TASK_SIZE")
	_ = viper.BindEnv(mongoConnStrKey, appEnvMongoPrefix+"_CONNSTR")
	_ = viper.BindEnv(mongoDbNameKey, appEnvMongoPrefix+"_DBNAME")
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
