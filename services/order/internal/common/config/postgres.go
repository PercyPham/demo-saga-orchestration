package config

import "services.order/internal/common/config/internal"

type PostgreSQLConfig struct {
	User     string `env:"POSTGRESQL_USER" env-default:"admin"`
	Password string `env:"POSTGRESQL_PASSWORD" env-default:"password"`
	DB       string `env:"POSTGRESQL_DB" env-default:"xemmenu_order_service_db"`
	Host     string `env:"POSTGRESQL_HOST" env-default:"localhost"`
	Port     string `env:"POSTGRESQL_PORT" env-default:"5432"`
}

var postgreSQLConfig PostgreSQLConfig

func PostgreSQL() PostgreSQLConfig {
	return postgreSQLConfig
}

func init() {
	internal.ReadEnv(&postgreSQLConfig)
}
