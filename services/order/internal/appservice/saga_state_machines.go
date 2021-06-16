package appservice

import (
	"github.com/percypham/saga-go"
	"services.order/internal/appservice/create_order"
)

var sagaStateMachines = []saga.StateMachine{
	create_order.NewCreateOrderStateMachine(),
}

func RegisterStateMachines(sagaManager saga.Manager) {
	for _, machine := range sagaStateMachines {
		sagaManager.Register(machine)
	}
}
