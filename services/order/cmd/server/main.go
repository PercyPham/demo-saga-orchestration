package main

import (
	"github.com/percypham/saga-go"
	"services.order/internal/adapter/db/postgresql"
	"services.order/internal/adapter/http/rest"
	"services.order/internal/adapter/mq"
	"services.order/internal/adapter/rabbitmq"
	"services.order/internal/common/config"
	"services.shared/logger/consolelogger"
)

func main() {
	log := consolelogger.New()

	repo, err := postgresql.Connect(config.PostgreSQL())
	if err != nil {
		log.Fatal("cannot connect to Postgres DB:", err)
		panic(err)
	}

	outflowConn, inflowConn, _, err := rabbitmq.Connect(config.RabbitMQ())
	if err != nil {
		log.Fatal("cannot connect to RabbitMQ:", err)
		panic(err)
	}

	sagaManager, err := saga.NewManager(saga.Config{
		SagaRepo:       repo,
		MQProducer:     mq.NewProducer(outflowConn),
		MQConsumer:     mq.NewConsumer(inflowConn),
		CommandChannel: config.Saga().CommandChannel,
		ReplyChannel:   config.Saga().ReplyChannel,
	})
	if err != nil {
		log.Fatal("cannot create sagaManager:", err)
		panic(err)
	}

	//sagaManager.RegisterCommandHandlers(handlers)
	//sagaManager.RegisterStateMachines(machines)

	go sagaManager.Serve()

	orderRestApiServer := rest.NewOrderRestApiServer(log, repo)

	err = orderRestApiServer.Run()
	if err != nil {
		log.Fatal("cannot run order rest api server:", err)
	}
}
