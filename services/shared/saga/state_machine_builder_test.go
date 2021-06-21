package saga_test

import (
	"services.shared/saga"
	"testing"
)

func TestStateMachineBuilder(t *testing.T) {
	stateMachine, _ := saga.StateMachineBuilder().
		For("CreateOrderSaga").
		WithCompensation(nil).
		InvokeParticipant(nil).
		OnReply("OutOfStock", nil).
		Build()

	if stateMachine.SagaType() != "CreateOrderSaga" {
		t.Errorf("expected CreateOrderSaga, got %v", stateMachine.SagaType())
	}
}
