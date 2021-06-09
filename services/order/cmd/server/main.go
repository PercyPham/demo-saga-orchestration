package main

import (
	"services.order/internal/adapter/db/postgresql"
	"services.order/internal/adapter/http/rest"
	"services.order/internal/common/config"
	"services.shared/logger/consolelogger"
)

func main() {
	log := consolelogger.New()

	repo, err := postgresql.Connect(config.PostgreSQL())
	if err != nil {
		log.Fatal("cannot connect to Postgres DB")
		panic(err)
	}

	rest.RunOrderServer(repo)
}
