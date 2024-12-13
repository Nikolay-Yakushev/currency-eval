package main

import (
	"fmt"
	"github.com/spf13/viper"
)

type Config struct {
	LogLevel              string `mapstructure:"LOG_LEVEL"`
	CurrencyServiceApiKey string `mapstructure:"CURRENCY_SERVICE_API_KEY"`
	PostgresHost          string `mapstructure:"POSTGRES_HOST"`
	PostgresUser          string `mapstructure:"POSTGRES_USER"`
	PostgresPassword      string `mapstructure:"POSTGRES_PASSWORD"`
	PostgresDB            string `mapstructure:"POSTGRES_DB"`
	PostgresPort          int    `mapstructure:"POSTGRES_PORT"`
	RestApiPort           int    `mapstructure:"REST_API_PORT"`
}

func NewConf(path string) (Config, error) {
	var (
		err    error
		config Config
	)
	viper.AddConfigPath(path)
	viper.SetConfigFile(".env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return Config{}, fmt.Errorf("failed to load config. Reason %w", err)
	}

	if err = viper.Unmarshal(&config); err != nil {
		return Config{}, fmt.Errorf("failed to unmarshal cofig. Reason %w", err)
	}
	return config, nil

}
