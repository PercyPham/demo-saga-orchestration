package config

import "services.order/internal/common/config/internal"

type AppConfig struct {
	ENV  string `env:"APP_ENV" env-default:"development"`
	PORT int    `env:"APP_PORT" env-default:"5000"`
}

var appConfig AppConfig

func App() AppConfig {
	return appConfig
}

func init() {
	internal.ReadEnv(&appConfig)
}
