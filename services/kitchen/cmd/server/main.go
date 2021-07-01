package main

import (
	"services.kitchen/internal/adapter/db/postgresql"
	"services.kitchen/internal/adapter/http/rest"
	"services.kitchen/internal/adapter/mq"
	"services.kitchen/internal/adapter/rabbitmq"
	"services.kitchen/internal/appservice"
	"services.kitchen/internal/common/config"
	"services.kitchen_contract"
	"services.shared/logger/consolelogger"
	"services.shared/saga"
)

func main() {
	log := consolelogger.New()

	repo, err := postgresql.Connect(config.PostgreSQL())
	if err != nil {
		log.Fatal("cannot connect to Postgres DB:", err)
		panic(err)
	}

	mqInflowConn, _, err := rabbitmq.Connect(config.RabbitMQ())
	if err != nil {
		panic("cannot connect MQ Inflow Connection: " + err.Error())
	}

	mqOutflowConn, _, err := rabbitmq.Connect(config.RabbitMQ())
	if err != nil {
		panic("cannot connect MQ Outflow Connection: " + err.Error())
	}

	sagaCommandHandler, err := saga.NewCommandHandler(saga.CommandHandlerConfig{
		CommandChannel: kitchen_contract.KitchenServiceCommandChannel,
		Producer:       mq.NewProducer(mqOutflowConn),
		Consumer:       mq.NewConsumer(mqInflowConn),
		MessageRepo:    repo,
	})
	if err != nil {
		log.Fatal("cannot create sagaCommandHandler:", err)
		panic(err)
	}

	appservice.HandleCommands(sagaCommandHandler, repo)

	go sagaCommandHandler.Serve()

	orderRestApiServer := rest.NewKitchenRestApiServer(log, repo, sagaCommandHandler)

	err = orderRestApiServer.Run()
	if err != nil {
		log.Fatal("cannot run order rest api server:", err)
	}
}
