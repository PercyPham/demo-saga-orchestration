package config

import "services.kitchen/internal/common/config/internal"

type RabbitMQConfig struct {
	User     string `env:"RABBITMQ_USER" env-default:"user"`
	Password string `env:"RABBITMQ_PASSWORD" env-default:"password"`
	Host     string `env:"RABBITMQ_HOST" env-default:"localhost"`
	Port     string `env:"RABBITMQ_PORT" env-default:"5672"`
}

var rabbitMQConfig RabbitMQConfig

func RabbitMQ() RabbitMQConfig {
	return rabbitMQConfig
}

func init() {
	internal.ReadEnv(&rabbitMQConfig)
}
