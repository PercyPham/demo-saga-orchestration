package config

import "services.order/internal/common/config/internal"

type AppConfig struct {
	ENV  string `env:"APP_ENV" env-default:"development"`
}

var appConfig AppConfig

func App() AppConfig {
	return appConfig
}

func init() {
	internal.ReadEnv(&appConfig)
}
