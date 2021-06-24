package main

import (
	"services.payment/internal/adapter/db/postgresql"
	"services.payment/internal/adapter/http/rest"
	"services.payment/internal/adapter/mq"
	"services.payment/internal/adapter/rabbitmq"
	"services.payment/internal/common/config"
	"services.payment_contract"
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
		CommandChannel: payment_contract.PaymentServiceCommandChannel,
		ReplyChannel:   payment_contract.PaymentServiceReplyChannel,
	})
	if err != nil {
		log.Fatal("cannot create sagaManager:", err)
		panic(err)
	}

	go sagaManager.Serve()

	paymentRestApiServer := rest.NewPaymentRestApiServer(log, repo, sagaManager)

	err = paymentRestApiServer.Run()
	if err != nil {
		log.Fatal("cannot run payment rest api server:", err)
	}
}
