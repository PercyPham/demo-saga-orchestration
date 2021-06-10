package config

import "services.order/internal/common/config/internal"

type SagaConfig struct {
	CommandChannel string `env:"SAGA_COMMAND_CHANNEL" env-default:"OrderService.SagaCommandChannel"`
	ReplyChannel   string `env:"SAGA_REPLY_CHANNEL" env-default:"OrderService.SagaReplyChannel"`
}

var sagaConfig SagaConfig

func Saga() SagaConfig {
	return sagaConfig
}

func init() {
	internal.ReadEnv(&sagaConfig)
}
