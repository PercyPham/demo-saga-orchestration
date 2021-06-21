package main

import (
	"services.order/internal/adapter/db/postgresql"
	"services.order/internal/adapter/http/rest"
	"services.order/internal/adapter/mq"
	"services.order/internal/adapter/rabbitmq"
	"services.order/internal/appservice"
	"services.order/internal/common/config"
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

	inflowConn, _, err := rabbitmq.Connect(config.RabbitMQ())
	if err != nil {
		panic("cannot connect MQ Inflow Connection: " + err.Error())
	}

	outflowConn, _, err := rabbitmq.Connect(config.RabbitMQ())
	if err != nil {
		panic("cannot connect MQ Outflow Connection: " + err.Error())
	}

	sagaManager, err := saga.NewManager(saga.Config{
		SagaRepo:       repo,
		Producer:       mq.NewProducer(outflowConn),
		Consumer:       mq.NewConsumer(inflowConn),
		CommandChannel: config.Saga().CommandChannel,
		ReplyChannel:   config.Saga().ReplyChannel,
	})
	if err != nil {
		log.Fatal("cannot create sagaManager:", err)
		panic(err)
	}

	appservice.RegisterStateMachines(sagaManager)

	go sagaManager.Serve()

	orderRestApiServer := rest.NewOrderRestApiServer(log, repo, sagaManager)

	err = orderRestApiServer.Run()
	if err != nil {
		log.Fatal("cannot run order rest api server:", err)
	}
}
