package config

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
)

func configureEnvs() {
	_ = viper.BindEnv("server.host", appEnvServerPrefix+"_HOST")
	_ = viper.BindEnv("server.port", appEnvServerPrefix+"_PORT")
	_ = viper.BindEnv("manager.host", appEnvManagerPrefix+"_HOST")
	_ = viper.BindEnv("manager.port", appEnvManagerPrefix+"_PORT")
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

func GetManagerHostAndPort() (string, error) {
	host := viper.GetString("manager.host")
	if host == "" {
		return "", ErrEmptyManagerHost
	}
	port := viper.GetString("manager.port")
	if port == "" {
		return "", ErrEmptyManagerPort
	}
	return fmt.Sprintf("%s:%s", host, port), nil
}
