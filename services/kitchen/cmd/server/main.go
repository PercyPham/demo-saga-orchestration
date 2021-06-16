package main

import (
	"github.com/percypham/saga-go"
	"services.kitchen/internal/adapter/http/rest"
	"services.kitchen/internal/common/config"
	"services.shared/logger/consolelogger"
)

func main() {

	log := consolelogger.New()

	_, err := saga.NewManager(saga.Config{
		SagaRepo:       nil,
		Producer:       nil,
		Consumer:       nil,
		CommandChannel: config.Saga().CommandChannel,
		ReplyChannel:   config.Saga().ReplyChannel,
	})
	//if err != nil {
	//	log.Fatal("cannot create sagaManager:", err)
	//	panic(err)
	//}


	orderRestApiServer := rest.NewKitchenRestApiServer(log)

	err = orderRestApiServer.Run()
	if err != nil {
		log.Fatal("cannot run order rest api server:", err)
	}
}