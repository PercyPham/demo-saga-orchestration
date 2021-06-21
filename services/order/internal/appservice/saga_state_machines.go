package appservice

import (
	"services.order/internal/appservice/create_order"
	"services.shared/saga"
)

var sagaStateMachines = []saga.StateMachine{
	create_order.NewCreateOrderStateMachine(),
}

func RegisterStateMachines(sagaManager saga.Manager) {
	for _, machine := range sagaStateMachines {
		sagaManager.Register(machine)
	}
}
