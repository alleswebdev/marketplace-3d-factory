package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Port         int
	Host         string
	Name         string
	WbToken      string
	OzonToken    string
	OzonClientID string
	DatabaseURL  string
}

func GetAppConfig() Config {
	cfg := Config{}

	viper.SetConfigName("values")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("../configs")
	viper.AddConfigPath("./configs")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	err = viper.Unmarshal(&cfg)
	if err != nil {
		panic(fmt.Errorf("fatal unmarshal config file: %w", err))
	}

	return cfg
}
