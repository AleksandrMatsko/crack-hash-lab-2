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
	_ = viper.BindEnv("server.host", appEnvServerPrefix+"_HOST")
	_ = viper.BindEnv("server.port", appEnvServerPrefix+"_PORT")
	_ = viper.BindEnv("workers.list", appEnvWorkersPrefix+"_LIST")
	_ = viper.BindEnv("workers.taskSize", appEnvWorkersPrefix+"_TASK_SIZE")
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
	s := viper.GetString("workers.list")
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
	s := viper.GetString("workers.taskSize")
	val, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0, err
	}
	return val, nil
}
