package config

import "services.payment/internal/common/config/internal"

type AppConfig struct {
	ENV  string `env:"APP_ENV" env-default:"development"`
	PORT int    `env:"APP_PORT" env-default:"5002"`
}

var appConfig AppConfig

func App() AppConfig {
	return appConfig
}

func init() {
	internal.ReadEnv(&appConfig)
}
