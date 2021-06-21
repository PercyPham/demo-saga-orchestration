package saga_test

import (
	"services.shared/saga"
	"services.shared/saga/msg"
	"testing"
)

func TestNewStateMachine(t *testing.T) {
	stateMachine, err := saga.StateMachineBuilder().For("SagaType").
		WithCompensation(func(sagaData []byte) (msg.Command, error) { return nil, nil }).
		InvokeParticipant(func(sagaData []byte) (msg.Command, error) { return nil, nil }).
		OnReply("replyType", func(sagaData []byte, reply msg.Reply) error { return nil }).
		Build()

	if err != nil {
		t.Errorf("expected no error, got error: %v", err)
	}

	if stateMachine.SagaType() != "SagaType" {
		t.Errorf("wrong saga type, expected SagaType, got %s", stateMachine.SagaType())
	}
}
