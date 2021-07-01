package main

import (
	"services.order/internal/adapter/db/postgresql"
	"services.order/internal/adapter/http/rest"
	"services.order/internal/adapter/mq"
	"services.order/internal/adapter/rabbitmq"
	"services.order/internal/appservice"
	"services.order/internal/common/config"
	"services.order_contract"
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

	sagaCommandHandler, err := saga.NewCommandHandler(saga.CommandHandlerConfig{
		CommandChannel: order_contract.OrderServiceCommandChannel,
		Producer:       mq.NewProducer(outflowConn),
		Consumer:       mq.NewConsumer(inflowConn),
		MessageRepo:    repo,
	})
	if err != nil {
		log.Fatal("cannot create sagaCommandHandler:", err)
		panic(err)
	}

	sagaManager, err := saga.NewManager(saga.Config{
		SagaRepo:       repo,
		Producer:       mq.NewProducer(outflowConn),
		Consumer:       mq.NewConsumer(inflowConn),
		ReplyChannel:   order_contract.OrderServiceReplyChannel,
	})
	if err != nil {
		log.Fatal("cannot create sagaManager:", err)
		panic(err)
	}

	appservice.HandleCommands(sagaCommandHandler, repo)
	appservice.RegisterStateMachines(sagaManager)

	go sagaCommandHandler.Serve()
	go sagaManager.Serve()

	orderRestApiServer := rest.NewOrderRestApiServer(log, repo, sagaManager)

	err = orderRestApiServer.Run()
	if err != nil {
		log.Fatal("cannot run order rest api server:", err)
	}
}
