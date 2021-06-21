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

	sagaManager, err := saga.NewManager(saga.Config{
		SagaRepo:       repo,
		Producer:       mq.NewProducer(mqOutflowConn),
		Consumer:       mq.NewConsumer(mqInflowConn),
		CommandChannel: kitchen_contract.KitchenServiceCommandChannel,
		ReplyChannel:   kitchen_contract.KitchenServiceReplyChannel,
	})
	if err != nil {
		log.Fatal("cannot create sagaManager:", err)
		panic(err)
	}

	appservice.HandleCommands(sagaManager, repo)

	go sagaManager.Serve()

	orderRestApiServer := rest.NewKitchenRestApiServer(log, repo)

	err = orderRestApiServer.Run()
	if err != nil {
		log.Fatal("cannot run order rest api server:", err)
	}
}
